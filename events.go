package ecs

import "strings"

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

func (t EventType) String() string {
	switch t {
	case EvtNone:
		return "EvtNone"
	case EvtComponentAdded:
		return "EvtComponentAdded"
	case EvtComponentRemoved:
		return "EvtComponentRemoved"
	case EvtComponentsResized:
		return "EvtComponentsResized"
	case EvtAny:
		return "EvtAny"
	}
	v := strings.Builder{}
	n := 0
	if t&EvtComponentAdded == EvtComponentAdded {
		if n > 0 {
			v.WriteRune('|')
		}
		v.WriteString("EvtComponentAdded")
		n++
	}
	if t&EvtComponentRemoved == EvtComponentRemoved {
		if n > 0 {
			v.WriteRune('|')
		}
		v.WriteString("EvtComponentRemoved")
		n++
	}
	if t&EvtComponentsResized == EvtComponentsResized {
		if n > 0 {
			v.WriteRune('|')
		}
		v.WriteString("EvtComponentsResized")
		n++
	}
	return v.String()
}
