package ecs

import (
	"reflect"
)

type eventStorage[T any] struct {
	frame uint64
	e0    []T
	e1    []T
}

func (es *eventStorage[T]) step() {
	var zt T
	es.frame++
	for len(es.e0) < len(es.e1) {
		es.e0 = append(es.e0, zt)
	}
	copy(es.e0, es.e1)
	es.e0 = es.e0[:len(es.e1)]
	es.e1 = es.e1[:0]
}

func (es *eventStorage[T]) add(t T) {
	es.e1 = append(es.e1, t)
}

func (es *eventStorage[T]) newReader() EventReaderFunc[T] {
	var zv T
	ecopy := make([]T, len(es.e0)+len(es.e1))
	copy(ecopy, es.e0)
	copy(ecopy[len(es.e0):], es.e1)
	return func() (T, bool) {
		if len(ecopy) == 0 {
			return zv, false
		}
		t := ecopy[0]
		ecopy = ecopy[1:]
		return t, true
	}
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
	ensureEventExists[T](ctx)
	var zt T
	evmap := ctx.world.getEvents()
	tm := typeMapKeyOf(reflect.TypeOf(zt))
	ew := evmap[tm].(*eventStorage[T])
	return func(t T) {
		ew.add(t)
	}
}

type EventWriterFunc[T any] func(T)

func EventReader[T any](ctx *Context) EventReaderFunc[T] {
	ensureEventExists[T](ctx)
	var zv T
	evmap := ctx.world.getEvents()

	tm := typeMapKeyOf(reflect.TypeOf(zv))
	ch := evmap[tm].(*eventStorage[T])
	return ch.newReader()
}

// EventReaderFunc returns false if there are no more events to read.
type EventReaderFunc[T any] func() (T, bool)
