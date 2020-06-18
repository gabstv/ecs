package ecs

type EventType uint

const (
	EvtNone              EventType = iota
	EvtComponentAdded    EventType = 1 << iota
	EvtComponentRemoved  EventType = 1 << iota
	EvtComponentsResized EventType = 1 << iota
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
