package ecs

import (
	"reflect"
)

type queryBase interface {
	Next() bool
	Reset()
	World() World
}

type Query1[T1 Component] interface {
	queryBase
	Clone() Query1[T1]
	Item() (Entity, *T1)
}

type Query2[T1, T2 Component] interface {
	queryBase
	Clone() Query2[T1, T2]
	Item() (Entity, *T1, *T2)
}

type Query3[T1, T2, T3 Component] interface {
	queryBase
	Clone() Query3[T1, T2, T3]
	Item() (Entity, *T1, *T2, *T3)
}

type Query4[T1, T2, T3, T4 Component] interface {
	queryBase
	Clone() Query4[T1, T2, T3, T4]
	Item() (Entity, *T1, *T2, *T3, *T4)
}

type Query5[T1, T2, T3, T4, T5 Component] interface {
	queryBase
	Clone() Query5[T1, T2, T3, T4, T5]
	Item() (Entity, *T1, *T2, *T3, *T4, *T5)
}

type worldQuery1[T1 Component] struct {
	w       World
	cs1     *componentStorage[T1]
	cursor1 int
}

func (wq *worldQuery1[T1]) Next() bool {
	wq.cursor1++
	if wq.cursor1 >= len(wq.cs1.Items) {
		return false
	}
	for wq.cs1.Items[wq.cursor1].IsDeleted {
		wq.cursor1++
		if wq.cursor1 >= len(wq.cs1.Items) {
			return false
		}
	}
	return true
}

func (wq *worldQuery1[T1]) Reset() {
	wq.cursor1 = -1
}

func (wq *worldQuery1[T1]) World() World {
	return wq.w
}

func (wq *worldQuery1[T1]) Clone() Query1[T1] {
	return &worldQuery1[T1]{
		w:       wq.w,
		cs1:     wq.cs1,
		cursor1: -1,
	}
}

func (wq *worldQuery1[T1]) Item() (Entity, *T1) {
	ref1 := &wq.cs1.Items[wq.cursor1]
	return ref1.Entity, &ref1.Component
}

// getOrCreateQuery1 also increases the usage by 1
func getOrCreateQuery1[T1 Component](w World) Query1[T1] {
	var zt T1
	rt := reflect.TypeOf(zt)
	th := typeTapeOf(rt)
	iq := w.getQuery(th)
	if iq != nil {
		q := iq.(Query1[T1])
		return q
	}
	// we also need to create the component storage if it doesn't exist yet
	// this is because systems can be registered before components are added
	// to the world
	// That is why we call getOrCreateComponentStorage

	cs1 := getOrCreateComponentStorage[T1](w)

	q := &worldQuery1[T1]{
		w:       w,
		cs1:     cs1,
		cursor1: -1,
	}

	return q
}

func Q1[T1 Component](ctx *Context) Query1[T1] {
	q := getOrCreateQuery1[T1](ctx.world)
	q.Reset()
	return q
}

// // //

type worldQuery2[T1, T2 Component] struct {
	w       World
	cs1     *componentStorage[T1]
	cs2     *componentStorage[T2]
	cursor1 int
	cursor2 int
}

func (wq *worldQuery2[T1, T2]) Next() bool {
beginNextWQ2:
	wq.cursor1++
	wq.cursor2++
	if wq.cursor1 >= len(wq.cs1.Items) || wq.cursor2 >= len(wq.cs2.Items) {
		return false
	}
	for wq.cs1.Items[wq.cursor1].Entity < wq.cs2.Items[wq.cursor2].Entity {
		wq.cursor1++
		if wq.cursor1 >= len(wq.cs1.Items) {
			return false
		}
	}
	for wq.cs2.Items[wq.cursor2].Entity < wq.cs1.Items[wq.cursor1].Entity {
		wq.cursor2++
		if wq.cursor2 >= len(wq.cs2.Items) {
			return false
		}
	}
	if wq.cs1.Items[wq.cursor1].IsDeleted || wq.cs2.Items[wq.cursor2].IsDeleted {
		goto beginNextWQ2
	}
	return true
}

func (wq *worldQuery2[T1, T2]) Reset() {
	wq.cursor1 = -1
	wq.cursor2 = -1
}

func (wq *worldQuery2[T1, T2]) World() World {
	return wq.w
}

func (wq *worldQuery2[T1, T2]) Clone() Query2[T1, T2] {
	return &worldQuery2[T1, T2]{
		w:       wq.w,
		cs1:     wq.cs1,
		cs2:     wq.cs2,
		cursor1: -1,
		cursor2: -1,
	}
}

func (wq *worldQuery2[T1, T2]) Item() (Entity, *T1, *T2) {
	ref1 := &wq.cs1.Items[wq.cursor1]
	ref2 := &wq.cs2.Items[wq.cursor2]
	return ref1.Entity, &ref1.Component, &ref2.Component
}

// getOrCreateQuery2 also increases the usage by 1
func getOrCreateQuery2[T1, T2 Component](w World) Query2[T1, T2] {
	var zt1 T1
	rt1 := reflect.TypeOf(zt1)
	var zt2 T2
	rt2 := reflect.TypeOf(zt2)
	th := typeTapeOf(rt1, rt2)
	iq := w.getQuery(th)
	if iq != nil {
		q := iq.(Query2[T1, T2])
		q.Reset()
		return q
	}
	// we also need to create the component storage if it doesn't exist yet
	// this is because systems can be registered before components are added
	// to the world
	// That is why we call getOrCreateComponentStorage

	cs1 := getOrCreateComponentStorage[T1](w)
	cs2 := getOrCreateComponentStorage[T2](w)

	q := &worldQuery2[T1, T2]{
		w:       w,
		cs1:     cs1,
		cs2:     cs2,
		cursor1: -1,
		cursor2: -1,
	}

	return q
}

func Q2[T1, T2 Component](ctx *Context) Query2[T1, T2] {
	q := getOrCreateQuery2[T1, T2](ctx.world)
	q.Reset()
	return q
}

// // //

type worldQuery3[T1, T2, T3 Component] struct {
	w       World
	cs1     *componentStorage[T1]
	cs2     *componentStorage[T2]
	cs3     *componentStorage[T3]
	cursor1 int
	cursor2 int
	cursor3 int
}

func (wq *worldQuery3[T1, T2, T3]) Next() bool {
beginNextWQ3:
	if wq.cursor1 >= len(wq.cs1.Items) || wq.cursor2 >= len(wq.cs2.Items) || wq.cursor3 >= len(wq.cs3.Items) {
		return false
	}
	wq.cursor1++
	wq.cursor2++
	wq.cursor3++
	for wq.cs1.Items[wq.cursor1].Entity < wq.cs2.Items[wq.cursor2].Entity {
		wq.cursor1++
		if wq.cursor1 >= len(wq.cs1.Items) {
			return false
		}
	}
	for wq.cs1.Items[wq.cursor1].Entity < wq.cs3.Items[wq.cursor3].Entity {
		wq.cursor1++
		if wq.cursor1 >= len(wq.cs1.Items) {
			return false
		}
	}
	for wq.cs2.Items[wq.cursor2].Entity < wq.cs1.Items[wq.cursor1].Entity {
		wq.cursor2++
		if wq.cursor2 >= len(wq.cs2.Items) {
			return false
		}
	}
	for wq.cs2.Items[wq.cursor2].Entity < wq.cs3.Items[wq.cursor3].Entity {
		wq.cursor2++
		if wq.cursor2 >= len(wq.cs2.Items) {
			return false
		}
	}
	for wq.cs3.Items[wq.cursor3].Entity < wq.cs1.Items[wq.cursor1].Entity {
		wq.cursor3++
		if wq.cursor3 >= len(wq.cs3.Items) {
			return false
		}
	}
	for wq.cs3.Items[wq.cursor3].Entity < wq.cs2.Items[wq.cursor2].Entity {
		wq.cursor3++
		if wq.cursor3 >= len(wq.cs3.Items) {
			return false
		}
	}
	if wq.cs1.Items[wq.cursor1].IsDeleted || wq.cs2.Items[wq.cursor2].IsDeleted || wq.cs3.Items[wq.cursor3].IsDeleted {
		goto beginNextWQ3
	}
	return true
}

func (wq *worldQuery3[T1, T2, T3]) Reset() {
	wq.cursor1 = -1
	wq.cursor2 = -1
	wq.cursor3 = -1
}

func (wq *worldQuery3[T1, T2, T3]) World() World {
	return wq.w
}

func (wq *worldQuery3[T1, T2, T3]) Clone() Query3[T1, T2, T3] {
	return &worldQuery3[T1, T2, T3]{
		w:       wq.w,
		cs1:     wq.cs1,
		cs2:     wq.cs2,
		cs3:     wq.cs3,
		cursor1: -1,
		cursor2: -1,
		cursor3: -1,
	}
}

func (wq *worldQuery3[T1, T2, T3]) Item() (Entity, *T1, *T2, *T3) {
	ref1 := &wq.cs1.Items[wq.cursor1]
	ref2 := &wq.cs2.Items[wq.cursor2]
	ref3 := &wq.cs3.Items[wq.cursor3]
	return ref1.Entity, &ref1.Component, &ref2.Component, &ref3.Component
}

// getOrCreateQuery3 also increases the usage by 1
func getOrCreateQuery3[T1, T2, T3 Component](w World) Query3[T1, T2, T3] {
	var zt1 T1
	rt1 := reflect.TypeOf(zt1)
	var zt2 T2
	rt2 := reflect.TypeOf(zt2)
	var zt3 T3
	rt3 := reflect.TypeOf(zt3)
	th := typeTapeOf(rt1, rt2, rt3)
	iq := w.getQuery(th)
	if iq != nil {
		q := iq.(Query3[T1, T2, T3])
		q.Reset()
		return q
	}
	// we also need to create the component storage if it doesn't exist yet
	// this is because systems can be registered before components are added
	// to the world
	// That is why we call getOrCreateComponentStorage

	cs1 := getOrCreateComponentStorage[T1](w)
	cs2 := getOrCreateComponentStorage[T2](w)
	cs3 := getOrCreateComponentStorage[T3](w)

	q := &worldQuery3[T1, T2, T3]{
		w:       w,
		cs1:     cs1,
		cs2:     cs2,
		cs3:     cs3,
		cursor1: -1,
		cursor2: -1,
		cursor3: -1,
	}

	return q
}

func Q3[T1, T2, T3 Component](ctx *Context) Query3[T1, T2, T3] {
	q := getOrCreateQuery3[T1, T2, T3](ctx.world)
	q.Reset()
	return q
}

// // //

type worldQuery4[T1, T2, T3, T4 Component] struct {
	w       World
	cs1     *componentStorage[T1]
	cs2     *componentStorage[T2]
	cs3     *componentStorage[T3]
	cs4     *componentStorage[T4]
	cursor1 int
	cursor2 int
	cursor3 int
	cursor4 int
}

func (wq *worldQuery4[T1, T2, T3, T4]) Next() bool {
beginNextWQ4:
	if wq.cursor1 >= len(wq.cs1.Items) || wq.cursor2 >= len(wq.cs2.Items) || wq.cursor3 >= len(wq.cs3.Items) || wq.cursor4 >= len(wq.cs4.Items) {
		return false
	}
	wq.cursor1++
	wq.cursor2++
	wq.cursor3++
	wq.cursor4++
	for wq.cs1.Items[wq.cursor1].Entity < wq.cs2.Items[wq.cursor2].Entity ||
		wq.cs1.Items[wq.cursor1].Entity < wq.cs3.Items[wq.cursor3].Entity ||
		wq.cs1.Items[wq.cursor1].Entity < wq.cs4.Items[wq.cursor4].Entity {
		wq.cursor1++
		if wq.cursor1 >= len(wq.cs1.Items) {
			return false
		}
	}
	for wq.cs2.Items[wq.cursor2].Entity < wq.cs1.Items[wq.cursor1].Entity ||
		wq.cs2.Items[wq.cursor2].Entity < wq.cs3.Items[wq.cursor3].Entity ||
		wq.cs2.Items[wq.cursor2].Entity < wq.cs4.Items[wq.cursor4].Entity {
		wq.cursor2++
		if wq.cursor2 >= len(wq.cs2.Items) {
			return false
		}
	}
	for wq.cs3.Items[wq.cursor3].Entity < wq.cs1.Items[wq.cursor1].Entity ||
		wq.cs3.Items[wq.cursor3].Entity < wq.cs2.Items[wq.cursor2].Entity ||
		wq.cs3.Items[wq.cursor3].Entity < wq.cs4.Items[wq.cursor4].Entity {
		wq.cursor3++
		if wq.cursor3 >= len(wq.cs3.Items) {
			return false
		}
	}
	for wq.cs4.Items[wq.cursor4].Entity < wq.cs1.Items[wq.cursor1].Entity ||
		wq.cs4.Items[wq.cursor4].Entity < wq.cs2.Items[wq.cursor2].Entity ||
		wq.cs4.Items[wq.cursor4].Entity < wq.cs3.Items[wq.cursor3].Entity {
		wq.cursor4++
		if wq.cursor4 >= len(wq.cs4.Items) {
			return false
		}
	}
	if wq.cs1.Items[wq.cursor1].IsDeleted || wq.cs2.Items[wq.cursor2].IsDeleted || wq.cs3.Items[wq.cursor3].IsDeleted || wq.cs4.Items[wq.cursor4].IsDeleted {
		goto beginNextWQ4
	}
	return true
}

func (wq *worldQuery4[T1, T2, T3, T4]) Reset() {
	wq.cursor1 = -1
	wq.cursor2 = -1
	wq.cursor3 = -1
	wq.cursor4 = -1
}

func (wq *worldQuery4[T1, T2, T3, T4]) World() World {
	return wq.w
}

func (wq *worldQuery4[T1, T2, T3, T4]) Clone() Query4[T1, T2, T3, T4] {
	return &worldQuery4[T1, T2, T3, T4]{
		w:       wq.w,
		cs1:     wq.cs1,
		cs2:     wq.cs2,
		cs3:     wq.cs3,
		cs4:     wq.cs4,
		cursor1: -1,
		cursor2: -1,
		cursor3: -1,
		cursor4: -1,
	}
}

func (wq *worldQuery4[T1, T2, T3, T4]) Item() (Entity, *T1, *T2, *T3, *T4) {
	ref1 := &wq.cs1.Items[wq.cursor1]
	ref2 := &wq.cs2.Items[wq.cursor2]
	ref3 := &wq.cs3.Items[wq.cursor3]
	ref4 := &wq.cs4.Items[wq.cursor4]
	return ref1.Entity, &ref1.Component, &ref2.Component, &ref3.Component, &ref4.Component
}

func getOrCreateQuery4[T1, T2, T3, T4 Component](w World) Query4[T1, T2, T3, T4] {
	var zt1 T1
	rt1 := reflect.TypeOf(zt1)
	var zt2 T2
	rt2 := reflect.TypeOf(zt2)
	var zt3 T3
	rt3 := reflect.TypeOf(zt3)
	var zt4 T4
	rt4 := reflect.TypeOf(zt4)
	th := typeTapeOf(rt1, rt2, rt3, rt4)
	iq := w.getQuery(th)
	if iq != nil {
		q := iq.(Query4[T1, T2, T3, T4])
		q.Reset()
		return q
	}
	// we also need to create the component storage if it doesn't exist yet
	// this is because systems can be registered before components are added
	// to the world
	// That is why we call getOrCreateComponentStorage

	cs1 := getOrCreateComponentStorage[T1](w)
	cs2 := getOrCreateComponentStorage[T2](w)
	cs3 := getOrCreateComponentStorage[T3](w)
	cs4 := getOrCreateComponentStorage[T4](w)

	q := &worldQuery4[T1, T2, T3, T4]{
		w:       w,
		cs1:     cs1,
		cs2:     cs2,
		cs3:     cs3,
		cs4:     cs4,
		cursor1: -1,
		cursor2: -1,
		cursor3: -1,
		cursor4: -1,
	}

	return q
}

func Q4[T1, T2, T3, T4 Component](ctx *Context) Query4[T1, T2, T3, T4] {
	q := getOrCreateQuery4[T1, T2, T3, T4](ctx.world)
	q.Reset()
	return q
}

// // //

type worldQuery5[T1, T2, T3, T4, T5 Component] struct {
	w       World
	cs1     *componentStorage[T1]
	cs2     *componentStorage[T2]
	cs3     *componentStorage[T3]
	cs4     *componentStorage[T4]
	cs5     *componentStorage[T5]
	cursor1 int
	cursor2 int
	cursor3 int
	cursor4 int
	cursor5 int
}

func (wq *worldQuery5[T1, T2, T3, T4, T5]) Next() bool {
beginNextWQ5:
	if wq.cursor1 >= len(wq.cs1.Items) || wq.cursor2 >= len(wq.cs2.Items) ||
		wq.cursor3 >= len(wq.cs3.Items) || wq.cursor4 >= len(wq.cs4.Items) || wq.cursor5 >= len(wq.cs5.Items) {
		return false
	}
	wq.cursor1++
	wq.cursor2++
	wq.cursor3++
	wq.cursor4++
	wq.cursor5++
	for wq.cs1.Items[wq.cursor1].Entity < wq.cs2.Items[wq.cursor2].Entity ||
		wq.cs1.Items[wq.cursor1].Entity < wq.cs3.Items[wq.cursor3].Entity ||
		wq.cs1.Items[wq.cursor1].Entity < wq.cs4.Items[wq.cursor4].Entity ||
		wq.cs1.Items[wq.cursor1].Entity < wq.cs5.Items[wq.cursor5].Entity {
		wq.cursor1++
		if wq.cursor1 >= len(wq.cs1.Items) {
			return false
		}
	}
	for wq.cs2.Items[wq.cursor2].Entity < wq.cs1.Items[wq.cursor1].Entity ||
		wq.cs2.Items[wq.cursor2].Entity < wq.cs3.Items[wq.cursor3].Entity ||
		wq.cs2.Items[wq.cursor2].Entity < wq.cs4.Items[wq.cursor4].Entity ||
		wq.cs2.Items[wq.cursor2].Entity < wq.cs5.Items[wq.cursor5].Entity {
		wq.cursor2++
		if wq.cursor2 >= len(wq.cs2.Items) {
			return false
		}
	}
	for wq.cs3.Items[wq.cursor3].Entity < wq.cs1.Items[wq.cursor1].Entity ||
		wq.cs3.Items[wq.cursor3].Entity < wq.cs2.Items[wq.cursor2].Entity ||
		wq.cs3.Items[wq.cursor3].Entity < wq.cs4.Items[wq.cursor4].Entity ||
		wq.cs3.Items[wq.cursor3].Entity < wq.cs5.Items[wq.cursor5].Entity {
		wq.cursor3++
		if wq.cursor3 >= len(wq.cs3.Items) {
			return false
		}
	}
	for wq.cs4.Items[wq.cursor4].Entity < wq.cs1.Items[wq.cursor1].Entity ||
		wq.cs4.Items[wq.cursor4].Entity < wq.cs2.Items[wq.cursor2].Entity ||
		wq.cs4.Items[wq.cursor4].Entity < wq.cs3.Items[wq.cursor3].Entity ||
		wq.cs4.Items[wq.cursor4].Entity < wq.cs5.Items[wq.cursor5].Entity {
		wq.cursor4++
		if wq.cursor4 >= len(wq.cs4.Items) {
			return false
		}
	}
	for wq.cs5.Items[wq.cursor5].Entity < wq.cs1.Items[wq.cursor1].Entity ||
		wq.cs5.Items[wq.cursor5].Entity < wq.cs2.Items[wq.cursor2].Entity ||
		wq.cs5.Items[wq.cursor5].Entity < wq.cs3.Items[wq.cursor3].Entity ||
		wq.cs5.Items[wq.cursor5].Entity < wq.cs4.Items[wq.cursor4].Entity {
		wq.cursor5++
		if wq.cursor5 >= len(wq.cs5.Items) {
			return false
		}
	}
	if wq.cs1.Items[wq.cursor1].IsDeleted || wq.cs2.Items[wq.cursor2].IsDeleted ||
		wq.cs3.Items[wq.cursor3].IsDeleted || wq.cs4.Items[wq.cursor4].IsDeleted ||
		wq.cs5.Items[wq.cursor5].IsDeleted {
		goto beginNextWQ5
	}
	return true
}

func (wq *worldQuery5[T1, T2, T3, T4, T5]) Reset() {
	wq.cursor1 = -1
	wq.cursor2 = -1
	wq.cursor3 = -1
	wq.cursor4 = -1
	wq.cursor5 = -1
}

func (wq *worldQuery5[T1, T2, T3, T4, T5]) World() World {
	return wq.w
}

func (wq *worldQuery5[T1, T2, T3, T4, T5]) Clone() Query5[T1, T2, T3, T4, T5] {
	return &worldQuery5[T1, T2, T3, T4, T5]{
		w:       wq.w,
		cs1:     wq.cs1,
		cs2:     wq.cs2,
		cs3:     wq.cs3,
		cs4:     wq.cs4,
		cs5:     wq.cs5,
		cursor1: -1,
		cursor2: -1,
		cursor3: -1,
		cursor4: -1,
		cursor5: -1,
	}
}

func (wq *worldQuery5[T1, T2, T3, T4, T5]) Item() (Entity, *T1, *T2, *T3, *T4, *T5) {
	ref1 := &wq.cs1.Items[wq.cursor1]
	ref2 := &wq.cs2.Items[wq.cursor2]
	ref3 := &wq.cs3.Items[wq.cursor3]
	ref4 := &wq.cs4.Items[wq.cursor4]
	ref5 := &wq.cs5.Items[wq.cursor5]
	return ref1.Entity, &ref1.Component, &ref2.Component, &ref3.Component, &ref4.Component, &ref5.Component
}

func getOrCreateQuery5[T1, T2, T3, T4, T5 Component](w World) Query5[T1, T2, T3, T4, T5] {
	var zt1 T1
	rt1 := reflect.TypeOf(zt1)
	var zt2 T2
	rt2 := reflect.TypeOf(zt2)
	var zt3 T3
	rt3 := reflect.TypeOf(zt3)
	var zt4 T4
	rt4 := reflect.TypeOf(zt4)
	var zt5 T5
	rt5 := reflect.TypeOf(zt5)
	th := typeTapeOf(rt1, rt2, rt3, rt4, rt5)
	iq := w.getQuery(th)
	if iq != nil {
		q := iq.(Query5[T1, T2, T3, T4, T5])
		q.Reset()
		return q
	}
	// we also need to create the component storage if it doesn't exist yet
	// this is because systems can be registered before components are added
	// to the world
	// That is why we call getOrCreateComponentStorage

	cs1 := getOrCreateComponentStorage[T1](w)
	cs2 := getOrCreateComponentStorage[T2](w)
	cs3 := getOrCreateComponentStorage[T3](w)
	cs4 := getOrCreateComponentStorage[T4](w)
	cs5 := getOrCreateComponentStorage[T5](w)

	q := &worldQuery5[T1, T2, T3, T4, T5]{
		w:       w,
		cs1:     cs1,
		cs2:     cs2,
		cs3:     cs3,
		cs4:     cs4,
		cs5:     cs5,
		cursor1: -1,
		cursor2: -1,
		cursor3: -1,
		cursor4: -1,
		cursor5: -1,
	}

	return q
}

func Q5[T1, T2, T3, T4, T5 Component](ctx *Context) Query5[T1, T2, T3, T4, T5] {
	q := getOrCreateQuery5[T1, T2, T3, T4, T5](ctx.world)
	q.Reset()
	return q
}
