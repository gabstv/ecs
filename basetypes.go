package ecs

type Entity uint64

type RegisterComponentFn func() BaseComponent
type RegisterSystemFn func() BaseSystem

type BaseComponent interface {
	UUID() string
	Name() string
	Flag() Flag
	Setup(w BaseWorld, f Flag, key [4]byte)
	Upsert(e Entity, data interface{})
	Remove(e Entity)
}

type BaseSystem interface {
	UUID() string
	Name() string
	ComponentAdded(e Entity, eflag Flag)
	ComponentRemoved(e Entity, eflag Flag)
	ComponentResized(cflag Flag)
	//V() View
	Priority() int64
	Setup(w BaseWorld)
	Enable()
	Disable()
	Enabled() bool
}

type BaseWorld interface {
	RegisterComponent(c BaseComponent)
	IsRegistered(id string) bool
	CFlag(e Entity) Flag
	NewEntity() Entity
	RemoveEntity(e Entity) bool
	C(id string) BaseComponent
	S(id string) BaseSystem
	CAdded(e Entity, c BaseComponent, key [4]byte)
	CRemoved(e Entity, c BaseComponent, key [4]byte)
	CResized(c BaseComponent, key [4]byte)
	AddSystem(s BaseSystem) error
	RemoveSystem(s BaseSystem)
	EachSystem(func(s BaseSystem) bool)
}

type View interface {
	//Matches() []Entity
}
