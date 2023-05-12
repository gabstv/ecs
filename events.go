package ecs

import (
	"fmt"
	"reflect"
)

var (
	DefaultEventQueueSize = 1024
)

// EnsureEventQueueSize ensures that the event queue size is of the given size.
// Make sure to call this function before any event writers are created.
func EnsureEventQueueSize[T any](ctx *Context, size int) {
	var zt T
	evmap := ctx.world.getEvents()
	tm := typeMapKeyOf(reflect.TypeOf(zt))
	evi := evmap[tm]
	if evi != nil {
		fmt.Printf("WARN: ecs.EnsureEventQueueSize[%T] called after one or more ecs.EventWriter[%T] calls. This may produce undesired results!\n", zt, zt)
		//TODO: check if we call close: close(evi.(chan T))
		// I opted to let the previous channel open to avoid a reader with an infinite loop to block forever.
	}
	ch := make(chan T, size)
	evmap[tm] = ch
}

func ensureEventExists[T any](ctx *Context) {
	var zt T
	evmap := ctx.world.getEvents()
	tm := typeMapKeyOf(reflect.TypeOf(zt))
	evi := evmap[tm]
	if evi == nil {
		ch := make(chan T, DefaultEventQueueSize)
		evmap[tm] = ch
	}
}

func EventWriter[T any](ctx *Context) EventWriterFunc[T] {
	ensureEventExists[T](ctx)
	var zt T
	evmap := ctx.world.getEvents()
	tm := typeMapKeyOf(reflect.TypeOf(zt))
	ch := evmap[tm].(chan T)
	return func(t T) bool {
		select {
		case ch <- t:
			return true
		default:
			return false
		}
	}
}

// EventWriterFunc will return false if the event was not written.
// The event may not be written if the event queue is full. To avoid this, make sure to allocate a
// big enough queue using ecs.EnsureEventQueueSize[T] (do this BEFORE creating any ecs.EventWriter[T]).
type EventWriterFunc[T any] func(T) bool

func EventReader[T any](ctx *Context) EventReaderFunc[T] {
	ensureEventExists[T](ctx)
	var zv T
	evmap := ctx.world.getEvents()

	tm := typeMapKeyOf(reflect.TypeOf(zv))
	ch := evmap[tm].(chan T)
	return func() (T, bool) {
		select {
		case t := <-ch:
			return t, true
		default:
			return zv, false
		}
	}
}

// EventReaderFunc returns false if there are no more events to read.
type EventReaderFunc[T any] func() (T, bool)
