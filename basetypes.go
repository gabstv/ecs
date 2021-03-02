package ecs

type Entity uint64

type RegisterComponentFn func() Component
type RegisterSystemFn func() System

type Component interface {
	UUID() string
	Name() string
	Flag() Flag
	Setup(w World, f Flag, key [4]byte)
	Upsert(e Entity, data interface{}) bool
	Remove(e Entity) bool
	World() World
}

// DispatchComponentEvent is a helper to dispatch component events
func DispatchComponentEvent(c Component, t EventType, e Entity) {
	c.World().Dispatch(Event{
		Type:          t,
		ComponentName: c.Name(),
		ComponentID:   c.UUID(),
		Entity:        e,
	})
}

type World interface {
	RegisterComponent(c Component)
	IsRegistered(id string) bool
	CFlag(e Entity) Flag
	NewEntity() Entity
	RemoveEntity(e Entity) bool
	C(id string) Component
	S(id string) System
	CAdded(e Entity, c Component, key [4]byte)
	CRemoved(e Entity, c Component, key [4]byte)
	CWillResize(c Component, key [4]byte)
	CResized(c Component, key [4]byte)
	AddSystem(s System) error
	RemoveSystem(s System)
	EachSystem(func(s System) bool)
	Dispatch(e Event)
	Listen(mask EventType, fn EventFn) int64
	RemoveListener(id int64)
	SetFlagGroup(name string, f Flag)
	FlagGroup(name string) Flag
	LGet(name string) interface{}
	LSet(name string, value interface{})
}

type View interface {
	//Matches() []Entity
}
