package ecs

type ComponentWatcher[T ComponentType] struct {
	ComponentAdded   func(e Entity)
	ComponentRemoved func(e Entity)

	comp *ComponentStore[T]
}

func newComponentWatcher[T ComponentType](comp *ComponentStore[T]) *ComponentWatcher[T] {
	w := &ComponentWatcher[T]{
		comp: comp,
	}
	comp.watchers.Add(w)
	return w
}

func (wa *ComponentWatcher[T]) Component() *ComponentStore[T] {
	return wa.comp
}

func (wa *ComponentWatcher[T]) Destroy() {
	if wa.comp == nil {
		panic("cannot destroy componentWatcher twice")
	}
	wa.comp.watchers.Remove(wa)
	wa.comp = nil
	wa.ComponentAdded = nil
	wa.ComponentRemoved = nil
}
