package ecs

import (
	"reflect"
	"sync"
)

// eventStorage may have a one frame delay if the reader system runs before the writer system
type eventStorage[T any] struct {
	lock  sync.RWMutex
	frame uint64
	e0    []eventStorageItem[T]
	e1    []eventStorageItem[T]
}

func (es *eventStorage[T]) step() {
	es.lock.Lock()
	defer es.lock.Unlock()
	var zt eventStorageItem[T]
	es.frame++
	for len(es.e0) < len(es.e1) {
		es.e0 = append(es.e0, zt)
	}
	copy(es.e0, es.e1)
	es.e0 = es.e0[:len(es.e1)]
	es.e1 = es.e1[:0]
}

func (es *eventStorage[T]) add(ctx *Context, t T) {
	es.lock.Lock()
	defer es.lock.Unlock()
	es.e1 = append(es.e1, eventStorageItem[T]{
		Data:        t,
		SystemIndex: ctx.currentSystemIndex,
	})
}

func (es *eventStorage[T]) newReader(ctx *Context) EventReaderFunc[T] {
	es.lock.RLock()
	defer es.lock.RUnlock()
	var zv T
	ecopy := make([]T, 0, len(es.e0)+len(es.e1))
	// es.e0 is always the previous frame, so we only add it if ctx.currentSystemIndex < es.e0.SystemIndex
	for _, v := range es.e0 {
		if ctx.currentSystemIndex < v.SystemIndex {
			ecopy = append(ecopy, v.Data)
		}
	}
	// es.e1 is always the current frame, so we only add it if ctx.currentSystemIndex >= es.e1.SystemIndex
	for _, v := range es.e1 {
		if ctx.currentSystemIndex >= v.SystemIndex {
			ecopy = append(ecopy, v.Data)
		}
	}
	return func() (T, bool) {
		if len(ecopy) == 0 {
			return zv, false
		}
		t := ecopy[0]
		ecopy = ecopy[1:]
		return t, true
	}
}

type eventStorageItem[T any] struct {
	Data        T
	SystemIndex int
}

type genericEventStorage interface {
	step()
}

func ensureEventExists[T any](ctx *Context) {
	var zt T
	evmap := ctx.world.getEvents()
	tm := typeMapKeyOf(reflect.TypeOf(zt))
	evi := evmap[tm]
	if evi == nil {
		evmap[tm] = &eventStorage[T]{}
	}
}

// EventWriter[T]
func EventWriter[T any](ctx *Context) EventWriterFunc[T] {
	ctx.eventRWLock.Lock()
	defer ctx.eventRWLock.Unlock()
	ensureEventExists[T](ctx)
	var zt T
	evmap := ctx.world.getEvents()
	tm := typeMapKeyOf(reflect.TypeOf(zt))
	ew := evmap[tm].(*eventStorage[T])
	return func(t T) {
		ew.add(ctx, t)
	}
}

type EventWriterFunc[T any] func(T)

func EventReader[T any](ctx *Context) EventReaderFunc[T] {
	ctx.eventRWLock.Lock()
	defer ctx.eventRWLock.Unlock()
	ensureEventExists[T](ctx)
	var zv T
	evmap := ctx.world.getEvents()

	tm := typeMapKeyOf(reflect.TypeOf(zv))
	ch := evmap[tm].(*eventStorage[T])
	return ch.newReader(ctx)
}

// EventReaderFunc returns false if there are no more events to read.
type EventReaderFunc[T any] func() (T, bool)

type EntityComponentPair[T any] struct {
	Entity Entity
	// This is a copy of the component at the time of the event.
	// If you absolutely need to get a pointer to the component, you can use

	ComponentCopy T
}

func getComponentAddedEventsParent[T Component](w World) *eventStorage[EntityComponentPair[T]] {
	var zt T
	zk := typeMapKeyOf(reflect.TypeOf(zt))
	m := w.getComponentAddedEvents()
	vi := m[zk]
	if vi == nil {
		vv := &eventStorage[EntityComponentPair[T]]{
			e0: make([]eventStorageItem[EntityComponentPair[T]], 0),
			e1: make([]eventStorageItem[EntityComponentPair[T]], 0),
		}
		vi = vv
		m[zk] = vi
	}
	v := vi.(*eventStorage[EntityComponentPair[T]])
	return v
}

func getComponentRemovedEventsParent[T Component](w World) *eventStorage[EntityComponentPair[T]] {
	var zt T
	zk := typeMapKeyOf(reflect.TypeOf(zt))
	m := w.getComponentRemovedEvents()
	vi := m[zk]
	if vi == nil {
		vv := &eventStorage[EntityComponentPair[T]]{
			e0: make([]eventStorageItem[EntityComponentPair[T]], 0),
			e1: make([]eventStorageItem[EntityComponentPair[T]], 0),
		}
		vi = vv
		m[zk] = vi
	}
	v := vi.(*eventStorage[EntityComponentPair[T]])
	return v
}

// ComponentsAdded returns a slice of the added components of the last frame.
func ComponentsAdded[T Component](ctx *Context) EventReaderFunc[EntityComponentPair[T]] {
	parent := getComponentAddedEventsParent[T](ctx.world)
	return parent.newReader(ctx)
}

// ComponentsRemoved returns a slice of the removed components of the last frame.
func ComponentsRemoved[T Component](ctx *Context) EventReaderFunc[EntityComponentPair[T]] {
	parent := getComponentRemovedEventsParent[T](ctx.world)
	return parent.newReader(ctx)
}
