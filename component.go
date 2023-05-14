package ecs

import (
	"reflect"

	"golang.org/x/exp/slices"
)

type Component interface {
	any
}

type ComponentWithCallback interface {
	OnAddedToEntity(e Entity)
}

type ComponentWithRemovalCallback interface {
	OnRemovedFromEntity(e Entity)
}

type componentStore[T Component] struct {
	Entity    Entity
	Component T
	IsDeleted bool
}

type componentStorage[T Component] struct {
	// needsSorting bool //TODO: remove this when proven unecessary
	zeroValue  T
	zeroType   reflect.Type
	mask       U256
	Items      []componentStore[T]
	references []*componentWeakReference[T]
}

func (s componentStorage[T]) ComponentType() reflect.Type {
	return s.zeroType
}

func (s componentStorage[T]) ComponentMask() U256 {
	return s.mask
}

func (s *componentStorage[T]) Add(e Entity, data T) {
	// call OnAddedToEntity on *T if it exists;
	// this is useful for components that need a reference of their own entity;
	if di, ok := any(&data).(ComponentWithCallback); ok {
		di.OnAddedToEntity(e)
	}

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
		// refresh weak references of entities that come after e
		if len(s.references) > 0 {
			for _, v := range s.references {
				if v.entity > e {
					v.refresh()
				}
			}
		}
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
		if di, ok := any(&(&s.Items[index]).Component).(ComponentWithRemovalCallback); ok {
			di.OnRemovedFromEntity(e)
		}
		// invalidate weak references!
		if len(s.references) > 0 {
			p, ok := slices.BinarySearchFunc(s.references, e, func(r *componentWeakReference[T], target Entity) int {
				if r.entity < target {
					return -1
				}
				if r.entity > target {
					return 1
				}
				return 0
			})
			if ok {
				s.references[p].lastIndex = -1
			}
		}
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

// getOrCreateWeakReference may returen nil if the entity does not have this component.
func (s *componentStorage[T]) getOrCreateWeakReference(e Entity) WeakReference[T] {
	if s.references == nil {
		s.references = make([]*componentWeakReference[T], 0)
	}
	index, ok := slices.BinarySearchFunc(s.references, e, func(r *componentWeakReference[T], target Entity) int {
		if r.entity < target {
			return -1
		}
		if r.entity > target {
			return 1
		}
		return 0
	})
	if ok {
		s.references[index].countIncr()
		return &componentWeakReferenceChild[T]{
			parent:  s.references[index],
			isValid: true,
		}
	}
	// create a new weak reference at index (if entity exists)
	entityIndex, _ := s.findEntity(e)
	if entityIndex == -1 {
		// not found!
		return nil
	}
	ref := &componentWeakReference[T]{
		parent:    s,
		entity:    e,
		lastIndex: entityIndex,
	}
	s.references = slices.Insert(s.references, index, ref)
	ref.countIncr()
	return &componentWeakReferenceChild[T]{
		parent:  ref,
		isValid: true,
	}
}

func (s *componentStorage[T]) removeWeakReference(r *componentWeakReference[T]) {
	if s.references == nil {
		panic("componentStorage[T].removeWeakReference called on nil references slice")
	}
	index, ok := slices.BinarySearchFunc(s.references, r.entity, func(r *componentWeakReference[T], target Entity) int {
		if r.entity < target {
			return -1
		}
		if r.entity > target {
			return 1
		}
		return 0
	})
	if !ok {
		panic("componentStorage[T].removeWeakReference called on invalid reference")
	}
	s.references = slices.Delete(s.references, index, index+1)
}

type componentWeakReference[T Component] struct {
	parent    *componentStorage[T]
	entity    Entity
	lastIndex int
	uses      int
}

func (r *componentWeakReference[T]) refresh() {
	index, _ := r.parent.findEntity(r.entity)
	r.lastIndex = index
}

func (r *componentWeakReference[T]) countIncr() {
	r.uses++
}

func (r *componentWeakReference[T]) countDecr() {
	r.uses--
	if r.uses == 0 {
		// remove from parent
		r.parent.removeWeakReference(r)
		r.parent = nil
		r.entity = 0
		r.lastIndex = -1
		return
	}
}

type componentWeakReferenceChild[T Component] struct {
	parent  *componentWeakReference[T]
	isValid bool
}

func (r *componentWeakReferenceChild[T]) Entity() Entity {
	if !r.isValid {
		panic("componentWeakReferenceChild[T].Entity() called on invalid reference")
	}
	return r.parent.entity
}

func (r *componentWeakReferenceChild[T]) Component() *T {
	if !r.isValid {
		panic("componentWeakReferenceChild[T].Component() called on invalid reference")
	}
	if r.parent.lastIndex == -1 {
		// entity was removed
		return nil
	}
	if r.parent.parent.Items[r.parent.lastIndex].Entity != r.parent.entity {
		//TODO: get the new index
		r.parent.lastIndex = -1
	}
	return &r.parent.parent.Items[r.parent.lastIndex].Component
}

func (r *componentWeakReferenceChild[T]) Destroy() {
	if !r.isValid {
		panic("componentWeakReferenceChild[T].Destroy() called on invalid reference")
	}
	r.parent.countDecr()
	r.parent = nil
	r.isValid = false
}

type WeakReference[T Component] interface {
	Entity() Entity
	Component() *T
	Destroy()
}

// NewWeakReference creates a weak reference to a component. This is useful when you need to
// create a node in a tree-like structure, so that you can easily find the parent node.
// Note that this is a relatively expensive operation when the component storagte needs to be sorted.
// It is recommended to use a query instead of this function when possible.
//
// After you obtain the reference, you NEED to call Destroy() when you are done with it (otherwise
//
//	the component storage will keep a pointer to this component reference forever).
func NewWeakReference[T Component](ctx *Context, entity Entity) WeakReference[T] {
	ct := getOrCreateComponentStorage[T](ctx.world)
	return ct.getOrCreateWeakReference(entity)
}
