package ecs

import (
	"reflect"

	"golang.org/x/exp/slices"
)

type Component interface {
	any
}

type componentStore[T Component] struct {
	Entity    Entity
	Component T
	IsDeleted bool
}

type componentStorage[T Component] struct {
	// needsSorting bool //TODO: remove this when proven unecessary
	zeroValue T
	zeroType  reflect.Type
	mask      U256
	Items     []componentStore[T]
}

func (s componentStorage[T]) ComponentType() reflect.Type {
	return s.zeroType
}

func (s componentStorage[T]) ComponentMask() U256 {
	return s.mask
}

func (s *componentStorage[T]) Add(e Entity, data T) {
	if len(s.Items) == 0 {
		s.Items = append(s.Items, componentStore[T]{
			Entity:    e,
			Component: data,
		})
		return
	}
	if s.Items[len(s.Items)-1].Entity < e {
		s.Items = append(s.Items, componentStore[T]{
			Entity:    e,
			Component: data,
		})
		return
	}
	if index, ok := slices.BinarySearchFunc(s.Items, e, func(item componentStore[T], e Entity) int {
		if item.Entity < e {
			return -1
		}
		if item.Entity > e {
			return 1
		}
		return 0
	}); ok {
		// replace existing componentStore[T] at index
		s.Items[index].Component = data
		s.Items[index].IsDeleted = false
	} else {
		// add a new componentStore[T] at index
		s.Items = slices.Insert(s.Items, index, componentStore[T]{
			Entity:    e,
			Component: data,
			IsDeleted: false,
		})
	}
}

func (s *componentStorage[T]) entityAt(index int) Entity {
	return s.Items[index].Entity
}

func (s *componentStorage[T]) removeEntity(e Entity) any {
	if index, ok := slices.BinarySearchFunc(s.Items, e, func(cs componentStore[T], target Entity) int {
		if cs.Entity < target {
			return -1
		}
		if cs.Entity > target {
			return 1
		}
		return 0
	}); ok {
		s.Items[index].IsDeleted = true
		return s.Items[index].Component
	}
	return nil
}

func (s *componentStorage[T]) fireComponentRemovedEvent(w World, e Entity, datacopy any) {
	getComponentRemovedEventsParent[T](w).add(EntityComponentPair[T]{
		Entity:        e,
		ComponentCopy: datacopy.(T),
	})
}

func (s *componentStorage[T]) findEntity(e Entity) (index int, data *T) {
	if index, ok := slices.BinarySearchFunc(s.Items, e, func(cs componentStore[T], target Entity) int {
		if cs.Entity < target {
			return -1
		}
		if cs.Entity > target {
			return 1
		}
		return 0
	}); ok {
		s.Items[index].IsDeleted = true
		if s.Items[index].IsDeleted {
			return -1, nil
		}
		return index, &(&s.Items[index]).Component
	}
	return -1, nil
}

type worldComponentStorage interface {
	ComponentType() reflect.Type
	ComponentMask() U256
	entityAt(index int) Entity
	removeEntity(e Entity) any
	// the only purpose of this function is to be called inside the world.removeEntity() function
	// This is because the component events are Generic, and the world cannot call generic methods.
	// The datacopy parameter is a struct copy of the component at the time of the event.
	fireComponentRemovedEvent(w World, e Entity, datacopy any)
}

func removeComponent[T Component](w World, e Entity) {
	cs := getOrCreateComponentStorage[T](w)
	fe := w.getFatEntity(e)
	if fe.ComponentMap.And(cs.ComponentMask()).IsZero() {
		return
	}
	fe.ComponentMap = fe.ComponentMap.AndNot(cs.ComponentMask())
	tv := cs.removeEntity(e)
	if tv == nil {
		return
	}
	d := tv.(T)
	getComponentRemovedEventsParent[T](w).add(EntityComponentPair[T]{
		Entity:        e,
		ComponentCopy: d,
	})
}

// GetComponent is a moderately expensive operation (in ECS terms) since it
// performs a binary search on the component storage. It is recommended to
// use a query instead of this function.
func GetComponent[T Component](ctx *Context, entity Entity) *T {
	ct := getOrCreateComponentStorage[T](ctx.world)
	_, ref := ct.findEntity(entity)
	return ref
}
