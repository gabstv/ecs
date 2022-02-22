package ecs

import "sync"

type Event struct {
	Name string
	Data interface{}
}

type ListenerID struct {
	Name string
	ID   int
}

type eventManager struct {
	l      sync.Mutex
	lastid int
	evts   map[string]map[int]func(e Event)
}

func newEventManager() *eventManager {
	return &eventManager{
		evts: make(map[string]map[int]func(e Event)),
	}
}

func (w *World) OnEvent(eventName string, fn func(e Event)) ListenerID {
	w.eventManager.l.Lock()
	defer w.eventManager.l.Unlock()
	if w.eventManager.evts[eventName] == nil {
		w.eventManager.evts[eventName] = make(map[int]func(e Event))
	}
	w.eventManager.lastid++
	id := w.eventManager.lastid
	w.eventManager.evts[eventName][id] = fn
	return ListenerID{
		Name: eventName,
		ID:   id,
	}
}

func (w *World) RemoveListener(id ListenerID) {
	w.eventManager.l.Lock()
	defer w.eventManager.l.Unlock()
	if w.eventManager.evts[id.Name] != nil {
		delete(w.eventManager.evts[id.Name], id.ID)
	}
}

func (w *World) FireEvent(eventName string, data interface{}) {
	w.eventManager.l.Lock()
	defer w.eventManager.l.Unlock()
	if w.eventManager.evts[eventName] != nil {
		for _, fn := range w.eventManager.evts[eventName] {
			fn(Event{
				Name: eventName,
				Data: data,
			})
		}
	}
}
