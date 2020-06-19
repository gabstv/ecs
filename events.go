package ecs

type EventType uint

const (
	EvtNone              EventType = 0
	EvtComponentAdded    EventType = 1 << 0
	EvtComponentRemoved  EventType = 1 << 1
	EvtComponentsResized EventType = 1 << 2
	EvtAny               EventType = 0b11111111111111111111111111111111
)

var allevents = [...]EventType{EvtComponentAdded, EvtComponentRemoved, EvtComponentsResized}

type Event struct {
	Type          EventType
	ComponentName string
	ComponentID   string
	Entity        Entity
}

type EventFn func(e Event)

type EventListener struct {
	ID   int64
	Mask EventType
	Fn   EventFn
}
