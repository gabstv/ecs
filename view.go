package ecs

type viewCommon struct {
	world    *World
	entities []Entity

	EntityAdded   func(e Entity)
	EntityRemoved func(e Entity)
}

func newViewCommon(w *World, onadded, onremoved func(e Entity)) *viewCommon {
	return &viewCommon{
		world:         w,
		entities:      make([]Entity, 0, 512),
		EntityAdded:   onadded,
		EntityRemoved: onremoved,
	}
}

func (vc *viewCommon) destroy() {
	vc.world = nil
	vc.entities = nil
	vc.EntityAdded = nil
	vc.EntityRemoved = nil
}

func (vc *viewCommon) World() *World {
	return vc.world
}

func (vc *viewCommon) Len() int {
	return len(vc.entities)
}

// ents return the internal array of entities. Used internally.
func (vc *viewCommon) ents() []Entity {
	return vc.entities
}

func (vc *viewCommon) onAdded(e Entity) {
	vc.EntityAdded(e)
}

func (vc *viewCommon) onRemoved(e Entity) {
	vc.EntityRemoved(e)
}

func (vc *viewCommon) removeEntityAt(index int) {
	vc.entities = append(vc.entities[:index], vc.entities[index+1:]...)
}

func (vc *viewCommon) addEntityAt(e Entity, index int) {
	Insert(vc.entities, index, e)
}

func (vc *viewCommon) entityIndex(e Entity) (int, bool) {
	return getEntityIndex(vc.entities, e)
}

type Viewer interface {
	World() *World
	Len() int

	ents() []Entity
	onAdded(e Entity)
	onRemoved(e Entity)
	removeEntityAt(index int)
	addEntityAt(e Entity, index int)
	entityIndex(e Entity) (int, bool)
}

type View[T ComponentType] struct {
	*viewCommon
	watcher *ComponentWatcher[T]
}

func (v *View[T]) Each(fn func(e Entity, d *T)) {
	// v.watcher.Component().Each(fn)
	slc := v.watcher.Component().all()
	for i := range slc {
		cd := &slc[i]
		fn(cd.Entity, &cd.Data)
	}
}

// Raw returns the raw component data for this view.
// WARNING: Use this only for fast copying the slice. Altering the original
// slice will most likely result in undefined behavior.
func (v *View[T]) Raw() []ComponentData[T] {
	return v.watcher.Component().all()
}

func (v *View[T]) Destroy() bool {
	if v.world == nil {
		return false
	}
	v.watcher.Destroy()
	v.watcher = nil
	v.viewCommon.destroy()
	return true
}

func buildWatcherAddedFunc(view Viewer, comps ...IComponentStore) func(e Entity) {
	return func(e Entity) {
		for _, c := range comps {
			if !c.Contains(e) {
				return
			}
		}
		if eindex, ok := view.entityIndex(e); !ok {
			view.addEntityAt(e, eindex)
			view.onAdded(e)
		}
	}
}

func buildWatcherRemovedFunc(view Viewer) func(e Entity) {
	return func(e Entity) {
		if eindex, ok := view.entityIndex(e); ok {
			view.removeEntityAt(eindex)
			view.onRemoved(e)
		}
	}
}

func NewView[T ComponentType](w *World, onadded, onremoved func(e Entity)) *View[T] {
	cc := GetComponentStore[T](w)
	view := &View[T]{
		viewCommon: newViewCommon(w, onadded, onremoved),
		watcher:    newComponentWatcher(cc),
	}
	view.watcher.ComponentAdded = buildWatcherAddedFunc(view)
	view.watcher.ComponentRemoved = buildWatcherRemovedFunc(view)
	// add all pre existing entities
	for _, cd := range cc.all() {
		view.entities = append(view.entities, cd.Entity)
		//TODO: maybe call ComponentAdded on all entities
	}
	return view
}

type View2[T1 ComponentType, T2 ComponentType] struct {
	*viewCommon
	watcher1 *ComponentWatcher[T1]
	watcher2 *ComponentWatcher[T2]
}

func NewView2[T1 ComponentType, T2 ComponentType](w *World, onadded, onremoved func(e Entity)) *View2[T1, T2] {
	cc1 := GetComponentStore[T1](w)
	cc2 := GetComponentStore[T2](w)
	view := &View2[T1, T2]{
		viewCommon: newViewCommon(w, onadded, onremoved),
		watcher1:   newComponentWatcher(cc1),
		watcher2:   newComponentWatcher(cc2),
	}
	view.watcher1.ComponentAdded = buildWatcherAddedFunc(view, cc1, cc2)
	view.watcher2.ComponentAdded = buildWatcherAddedFunc(view, cc1, cc2)
	view.watcher1.ComponentRemoved = buildWatcherRemovedFunc(view)
	view.watcher2.ComponentRemoved = buildWatcherRemovedFunc(view)
	// add all pre existing entities
	i1 := 0
	i2 := 0
	all1 := cc1.all()
	all2 := cc2.all()
	for {
		if i1 >= len(all1) || i2 >= len(all2) {
			break
		}
		if all1[i1].Entity < all2[i2].Entity {
			i1++
			continue
		}
		if all1[i1].Entity > all2[i2].Entity {
			i2++
			continue
		}
		if all1[i1].Entity == all2[i2].Entity {
			view.entities = append(view.entities, all1[i1].Entity)
		}
		i1++
		i2++
	}
	return view
}

func (v *View2[T1, T2]) Destroy() bool {
	if v.world == nil {
		return false
	}
	// close(v.closed)
	v.watcher1.Destroy()
	v.watcher2.Destroy()
	v.watcher1 = nil
	v.watcher2 = nil
	v.viewCommon.destroy()
	return true
}

func (v *View2[T1, T2]) Each(fn func(e Entity, d1 *T1, d2 *T2)) {
	c1 := v.watcher1.Component()
	c2 := v.watcher2.Component()
	ld1 := c1.all()
	len1 := len(ld1)
	ld2 := c2.all()
	len2 := len(ld2)
	i1 := 0
	i2 := 0

	for {
		if i1 >= len1 || i2 >= len2 {
			return
		}
		if ld1[i1].Entity == ld2[i2].Entity {
			cd1 := &ld1[i1]
			cd2 := &ld2[i2]
			fn(ld1[i1].Entity, &cd1.Data, &cd2.Data)
			i1++
			i2++
			continue
		}
		if ld1[i1].Entity < ld2[i2].Entity {
			i1++
		} else {
			i2++
		}
	}
}

type View3[T1 ComponentType, T2 ComponentType, T3 ComponentType] struct {
	*viewCommon
	watcher1 *ComponentWatcher[T1]
	watcher2 *ComponentWatcher[T2]
	watcher3 *ComponentWatcher[T3]
}

func NewView3[T1 ComponentType, T2 ComponentType, T3 ComponentType](w *World, onadded, onremoved func(e Entity)) *View3[T1, T2, T3] {
	cc1 := GetComponentStore[T1](w)
	cc2 := GetComponentStore[T2](w)
	cc3 := GetComponentStore[T3](w)
	view := &View3[T1, T2, T3]{
		viewCommon: newViewCommon(w, onadded, onremoved),
		watcher1:   newComponentWatcher(cc1),
		watcher2:   newComponentWatcher(cc2),
		watcher3:   newComponentWatcher(cc3),
	}
	view.watcher1.ComponentAdded = buildWatcherAddedFunc(view, cc1, cc2, cc3)
	view.watcher2.ComponentAdded = buildWatcherAddedFunc(view, cc1, cc2, cc3)
	view.watcher3.ComponentAdded = buildWatcherAddedFunc(view, cc1, cc2, cc3)
	view.watcher1.ComponentRemoved = buildWatcherRemovedFunc(view)
	view.watcher2.ComponentRemoved = buildWatcherRemovedFunc(view)
	view.watcher3.ComponentRemoved = buildWatcherRemovedFunc(view)
	// add all pre existing entities
	all1 := cc1.all()
	all2 := cc2.all()
	all3 := cc3.all()
	i1 := 0
	i2 := 0
	i3 := 0
	for {
		if i1 >= len(all1) || i2 >= len(all2) || i3 >= len(all3) {
			break
		}
		eres := ent3(all1[i1].Entity, all2[i2].Entity, all3[i3].Entity)
		switch eres {
		case 0:
			view.entities = append(view.entities, all1[i1].Entity)
			i1++
			i2++
			i3++
		case 1:
			i1++
		case 2:
			i2++
		case 3:
			i3++
		}
	}
	return view
}

func (v *View3[T1, T2, T3]) Destroy() bool {
	if v.world == nil {
		return false
	}
	v.watcher1.Destroy()
	v.watcher2.Destroy()
	v.watcher3.Destroy()
	v.watcher1 = nil
	v.watcher2 = nil
	v.watcher3 = nil
	v.viewCommon.destroy()
	return true
}

func (v *View3[T1, T2, T3]) Each(fn func(e Entity, d1 *T1, d2 *T2, d3 *T3)) {
	c1 := v.watcher1.Component()
	c2 := v.watcher2.Component()
	c3 := v.watcher3.Component()
	ld1 := c1.all()
	len1 := len(ld1)
	ld2 := c2.all()
	len2 := len(ld2)
	ld3 := c3.all()
	len3 := len(ld3)
	i1 := 0
	i2 := 0
	i3 := 0

	for {
		if i1 >= len1 || i2 >= len2 || i3 >= len3 {
			return
		}
		eres := ent3(ld1[i1].Entity, ld2[i2].Entity, ld3[i3].Entity)
		switch eres {
		case 0:
			cd1 := &ld1[i1]
			cd2 := &ld2[i2]
			cd3 := &ld3[i3]
			fn(ld1[i1].Entity, &cd1.Data, &cd2.Data, &cd3.Data)
			i1++
			i2++
			i3++
		case 1:
			i1++
		case 2:
			i2++
		case 3:
			i3++
		}
	}
}

type View4[T1 ComponentType, T2 ComponentType, T3 ComponentType, T4 ComponentType] struct {
	*viewCommon
	watcher1 *ComponentWatcher[T1]
	watcher2 *ComponentWatcher[T2]
	watcher3 *ComponentWatcher[T3]
	watcher4 *ComponentWatcher[T4]
}

func NewView4[T1 ComponentType, T2 ComponentType, T3 ComponentType, T4 ComponentType](w *World, onadded, onremoved func(e Entity)) *View4[T1, T2, T3, T4] {
	cc1 := GetComponentStore[T1](w)
	cc2 := GetComponentStore[T2](w)
	cc3 := GetComponentStore[T3](w)
	cc4 := GetComponentStore[T4](w)
	view := &View4[T1, T2, T3, T4]{
		viewCommon: newViewCommon(w, onadded, onremoved),
		watcher1:   newComponentWatcher(cc1),
		watcher2:   newComponentWatcher(cc2),
		watcher3:   newComponentWatcher(cc3),
		watcher4:   newComponentWatcher(cc4),
	}
	view.watcher1.ComponentAdded = buildWatcherAddedFunc(view, cc1, cc2, cc3, cc4)
	view.watcher2.ComponentAdded = buildWatcherAddedFunc(view, cc1, cc2, cc3, cc4)
	view.watcher3.ComponentAdded = buildWatcherAddedFunc(view, cc1, cc2, cc3, cc4)
	view.watcher4.ComponentAdded = buildWatcherAddedFunc(view, cc1, cc2, cc3, cc4)
	view.watcher1.ComponentRemoved = buildWatcherRemovedFunc(view)
	view.watcher2.ComponentRemoved = buildWatcherRemovedFunc(view)
	view.watcher3.ComponentRemoved = buildWatcherRemovedFunc(view)
	view.watcher4.ComponentRemoved = buildWatcherRemovedFunc(view)
	// add all pre existing entities
	all1 := cc1.all()
	all2 := cc2.all()
	all3 := cc3.all()
	all4 := cc4.all()
	i1 := 0
	i2 := 0
	i3 := 0
	i4 := 0
	for {
		if i1 >= len(all1) || i2 >= len(all2) || i3 >= len(all3) || i4 >= len(all4) {
			break
		}
		eres := ent4(all1[i1].Entity, all2[i2].Entity, all3[i3].Entity, all4[i4].Entity)
		switch eres {
		case 0:
			view.entities = append(view.entities, all1[i1].Entity)
			i1++
			i2++
			i3++
			i4++
		case 1:
			i1++
		case 2:
			i2++
		case 3:
			i3++
		case 4:
			i4++
		}
	}
	return view
}

func (v *View4[T1, T2, T3, T4]) Destroy() bool {
	if v.world == nil {
		return false
	}
	v.watcher1.Destroy()
	v.watcher2.Destroy()
	v.watcher3.Destroy()
	v.watcher4.Destroy()
	v.watcher1 = nil
	v.watcher2 = nil
	v.watcher3 = nil
	v.watcher4 = nil
	v.viewCommon.destroy()
	return true
}

func (v *View4[T1, T2, T3, T4]) Each(fn func(e Entity, d1 *T1, d2 *T2, d3 *T3, d4 *T4)) {
	c1 := v.watcher1.Component()
	c2 := v.watcher2.Component()
	c3 := v.watcher3.Component()
	c4 := v.watcher4.Component()
	ld1 := c1.all()
	len1 := len(ld1)
	ld2 := c2.all()
	len2 := len(ld2)
	ld3 := c3.all()
	len3 := len(ld3)
	ld4 := c4.all()
	len4 := len(ld4)
	i1 := 0
	i2 := 0
	i3 := 0
	i4 := 0

	for {
		if i1 >= len1 || i2 >= len2 || i3 >= len3 || i4 >= len4 {
			return
		}
		eres := ent4(ld1[i1].Entity, ld2[i2].Entity, ld3[i3].Entity, ld4[i4].Entity)
		switch eres {
		case 0:
			cd1 := &ld1[i1]
			cd2 := &ld2[i2]
			cd3 := &ld3[i3]
			cd4 := &ld4[i4]
			fn(ld1[i1].Entity, &cd1.Data, &cd2.Data, &cd3.Data, &cd4.Data)
			i1++
			i2++
			i3++
			i4++
		case 1:
			i1++
		case 2:
			i2++
		case 3:
			i3++
		case 4:
			i4++
		}
	}
}
