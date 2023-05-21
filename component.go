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
	references []*componentWeakReference[T]
	siblings   []worldComponentStorage
	items      []componentStore[T]
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

	if len(s.items) == 0 {
		if cap(s.items) == len(s.items) {
			// realloc will happen;
			// the func below will ensure that components that are often used together in the same memory region
			// have to play nice with the CPU cache
			s.ensureCapacity(len(s.items)*2, true)
		}
		s.items = append(s.items, componentStore[T]{
			Entity:    e,
			Component: data,
		})
		return
	}
	if s.items[len(s.items)-1].Entity < e {
		s.items = append(s.items, componentStore[T]{
			Entity:    e,
			Component: data,
		})
		return
	}
	if index, ok := slices.BinarySearchFunc(s.items, e, func(item componentStore[T], e Entity) int {
		if item.Entity < e {
			return -1
		}
		if item.Entity > e {
			return 1
		}
		return 0
	}); ok {
		// replace existing componentStore[T] at index
		s.items[index].Component = data
		s.items[index].IsDeleted = false
	} else {
		// add a new componentStore[T] at index
		s.items = slices.Insert(s.items, index, componentStore[T]{
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

func (s *componentStorage[T]) ensureCapacity(capacity int, forceRealloc bool) {
	if cap(s.items) >= capacity && !forceRealloc {
		return
	}
	//TODO: track the amount of reallocations that occured
	oldItems := s.items
	s.items = make([]componentStore[T], len(oldItems), max(capacity, cap(oldItems)*2))
	copy(s.items, oldItems)
	oldItems = nil
	for _, v := range s.siblings {
		v.reallocSelf(capacity)
	}
}

// this will only affect this component (no siblings)
func (s *componentStorage[T]) reallocSelf(capacity int) {
	oldItems := s.items
	s.items = make([]componentStore[T], len(oldItems), max(capacity, cap(oldItems)))
	copy(s.items, oldItems)
	oldItems = nil
}

func (s *componentStorage[T]) entityAt(index int) Entity {
	return s.items[index].Entity
}

func (s *componentStorage[T]) removeEntity(ctx *Context, e Entity) any {
	if index, ok := slices.BinarySearchFunc(s.items, e, func(cs componentStore[T], target Entity) int {
		if cs.Entity < target {
			return -1
		}
		if cs.Entity > target {
			return 1
		}
		return 0
	}); ok {
		s.items[index].IsDeleted = true
		if di, ok := any(&(&s.items[index]).Component).(ComponentWithRemovalCallback); ok {
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
		return s.items[index].Component
	}
	return nil
}

func (s *componentStorage[T]) fireComponentRemovedEvent(ctx *Context, e Entity, datacopy any) {
	w := ctx.world
	getComponentRemovedEventsParent[T](w).add(ctx, EntityComponentPair[T]{
		Entity:        e,
		ComponentCopy: datacopy.(T),
	})
}

func (s *componentStorage[T]) findEntity(e Entity) (index int, data *T) {
	if index, ok := slices.BinarySearchFunc(s.items, e, func(cs componentStore[T], target Entity) int {
		if cs.Entity < target {
			return -1
		}
		if cs.Entity > target {
			return 1
		}
		return 0
	}); ok {
		if s.items[index].IsDeleted {
			return -1, nil
		}
		return index, &(&s.items[index]).Component
	}
	return -1, nil
}

func (s *componentStorage[T]) gc() {
	icopy := make([]componentStore[T], 0, cap(s.items))
	for _, v := range s.items {
		if !v.IsDeleted {
			icopy = append(icopy, v)
		}
	}
	s.items = icopy
	for _, wr := range s.references {
		wr.refresh()
	}
}

func (s *componentStorage[T]) appendSibling(wc worldComponentStorage) {
	for _, v := range s.siblings {
		if v == wc {
			return
		}
	}
	s.siblings = append(s.siblings, wc)
}

type worldComponentStorage interface {
	ComponentType() reflect.Type
	ComponentMask() U256
	entityAt(index int) Entity
	ensureCapacity(capacity int, forceRealloc bool)
	reallocSelf(capacity int)
	removeEntity(ctx *Context, e Entity) any
	// the only purpose of this function is to be called inside the world.removeEntity() function
	// This is because the component events are Generic, and the world cannot call generic methods.
	// The datacopy parameter is a struct copy of the component at the time of the event.
	fireComponentRemovedEvent(ctx *Context, e Entity, datacopy any)
	gc()
}

func removeComponent[T Component](ctx *Context, e Entity) {
	w := ctx.world
	cs := getOrCreateComponentStorage[T](w, 0)
	fe := w.getFatEntity(e)
	if fe.ComponentMap.And(cs.ComponentMask()).IsZero() {
		return
	}
	fe.ComponentMap = fe.ComponentMap.AndNot(cs.ComponentMask())
	tv := cs.removeEntity(ctx, e)
	if tv == nil {
		return
	}
	d := tv.(T)
	getComponentRemovedEventsParent[T](w).add(ctx, EntityComponentPair[T]{
		Entity:        e,
		ComponentCopy: d,
	})
}

// EnsureComponentAffinity ensures that the component storage for this component is in the same memory region as the other components.
// This is useful to play nice with the CPU cache.
func EnsureComponentAffinity[T1, T2 Component](w World) {
	ct1 := getOrCreateComponentStorage[T1](w, 0)
	ct2 := getOrCreateComponentStorage[T2](w, 0)
	ct1.appendSibling(ct2)
	ct2.appendSibling(ct1)
}

// EnsureComponentCapacity must run before any AddComponent call to be effective.
// Use this when you have an idea of the average number of entities that will have this component.
// A smart use of this function can reduce the number of allocations.
func EnsureComponentCapacity[T Component](w World, capacity int) {
	ct := getOrCreateComponentStorage[T](w, capacity)
	ct.ensureCapacity(capacity, false)
}

// GetComponent is a moderately expensive operation (in ECS terms) since it
// performs a binary search on the component storage. It is recommended to
// use a query instead of this function.
func GetComponent[T Component](ctx *Context, entity Entity) *T {
	ct := getOrCreateComponentStorage[T](ctx.world, 0)
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
	if r.parent.parent.items[r.parent.lastIndex].Entity != r.parent.entity {
		//TODO: get the new index
		r.parent.lastIndex = -1
	}
	return &r.parent.parent.items[r.parent.lastIndex].Component
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
	ct := getOrCreateComponentStorage[T](ctx.world, 0)
	return ct.getOrCreateWeakReference(entity)
}
