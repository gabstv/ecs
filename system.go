package ecs

type SystemID uint64

type System func(*Context)

// AddStartupSystem adds a system that will be run once at the start of the
// nest Step() call.
// AddStartupSystem is thread safe.
func AddStartupSystem(w World, system func(*Context)) {
	w.addStartupSystem(system)
}

func AddSystem(w World, system func(*Context), opts ...AddSystemOptions) (SystemID, error) {
	o := addSystemOptions{}
	for _, v := range opts {
		v(&o)
	}
	return w.addSystem(worldSystem{
		SortPriority:   o.SortPriority,
		Value:          system,
		LocalResources: make(map[TypeMapKey]any),
	})
}

type addSystemOptions struct {
	SortPriority int
}

type AddSystemOptions func(*addSystemOptions)

func WithSortPriority(priority int) AddSystemOptions {
	return func(o *addSystemOptions) {
		o.SortPriority = priority
	}
}

type worldSystem struct {
	ID             SystemID
	SortPriority   int
	Value          System
	LocalResources map[TypeMapKey]any
}
