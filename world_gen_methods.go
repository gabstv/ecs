package ecs

func addComponent[T Component](ctx *Context, e Entity, data T) {
	w := ctx.world
	fatEntity := w.getFatEntity(e)
	assert(fatEntity != nil, "entity not found")
	cs := getOrCreateComponentStorage[T](w, 0)
	cs.Add(e, data)
	fatEntity.ComponentMap = fatEntity.ComponentMap.Or(cs.mask)
	// trigger events
	getComponentAddedEventsParent[T](w).add(ctx, EntityComponentPair[T]{
		Entity:        e,
		ComponentCopy: data,
	})
}

func addComponent2[T1, T2 Component](ctx *Context, e Entity, data1 T1, data2 T2) {
	w := ctx.world
	fatEntity := w.getFatEntity(e)
	assert(fatEntity != nil, "entity not found")
	cs1 := getOrCreateComponentStorage[T1](w, 0)
	cs2 := getOrCreateComponentStorage[T2](w, 0)
	cs1.Add(e, data1)
	cs2.Add(e, data2)
	fatEntity.ComponentMap = fatEntity.ComponentMap.Or(cs1.mask).Or(cs2.mask)
	// trigger events
	getComponentAddedEventsParent[T1](w).add(ctx, EntityComponentPair[T1]{
		Entity:        e,
		ComponentCopy: data1,
	})
	getComponentAddedEventsParent[T2](w).add(ctx, EntityComponentPair[T2]{
		Entity:        e,
		ComponentCopy: data2,
	})
}

func addComponent3[T1, T2, T3 Component](ctx *Context, e Entity, data1 T1, data2 T2, data3 T3) {
	w := ctx.world
	fatEntity := w.getFatEntity(e)
	assert(fatEntity != nil, "entity not found")
	cs1 := getOrCreateComponentStorage[T1](w, 0)
	cs2 := getOrCreateComponentStorage[T2](w, 0)
	cs3 := getOrCreateComponentStorage[T3](w, 0)
	cs1.Add(e, data1)
	cs2.Add(e, data2)
	cs3.Add(e, data3)
	fatEntity.ComponentMap = fatEntity.ComponentMap.Or(cs1.mask).Or(cs2.mask).Or(cs3.mask)
	// trigger events
	getComponentAddedEventsParent[T1](w).add(ctx, EntityComponentPair[T1]{
		Entity:        e,
		ComponentCopy: data1,
	})
	getComponentAddedEventsParent[T2](w).add(ctx, EntityComponentPair[T2]{
		Entity:        e,
		ComponentCopy: data2,
	})
	getComponentAddedEventsParent[T3](w).add(ctx, EntityComponentPair[T3]{
		Entity:        e,
		ComponentCopy: data3,
	})
}

func addComponent4[T1, T2, T3, T4 Component](ctx *Context, e Entity, data1 T1, data2 T2, data3 T3, data4 T4) {
	w := ctx.world
	fatEntity := w.getFatEntity(e)
	assert(fatEntity != nil, "entity not found")
	cs1 := getOrCreateComponentStorage[T1](w, 0)
	cs2 := getOrCreateComponentStorage[T2](w, 0)
	cs3 := getOrCreateComponentStorage[T3](w, 0)
	cs4 := getOrCreateComponentStorage[T4](w, 0)
	cs1.Add(e, data1)
	cs2.Add(e, data2)
	cs3.Add(e, data3)
	cs4.Add(e, data4)
	fatEntity.ComponentMap = fatEntity.ComponentMap.Or(cs1.mask).Or(cs2.mask).Or(cs3.mask).Or(cs4.mask)
	// trigger events
	getComponentAddedEventsParent[T1](w).add(ctx, EntityComponentPair[T1]{
		Entity:        e,
		ComponentCopy: data1,
	})
	getComponentAddedEventsParent[T2](w).add(ctx, EntityComponentPair[T2]{
		Entity:        e,
		ComponentCopy: data2,
	})
	getComponentAddedEventsParent[T3](w).add(ctx, EntityComponentPair[T3]{
		Entity:        e,
		ComponentCopy: data3,
	})
	getComponentAddedEventsParent[T4](w).add(ctx, EntityComponentPair[T4]{
		Entity:        e,
		ComponentCopy: data4,
	})
}

func addComponent5[T1, T2, T3, T4, T5 Component](ctx *Context, e Entity, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5) {
	w := ctx.world
	fatEntity := w.getFatEntity(e)
	assert(fatEntity != nil, "entity not found")
	cs1 := getOrCreateComponentStorage[T1](w, 0)
	cs2 := getOrCreateComponentStorage[T2](w, 0)
	cs3 := getOrCreateComponentStorage[T3](w, 0)
	cs4 := getOrCreateComponentStorage[T4](w, 0)
	cs5 := getOrCreateComponentStorage[T5](w, 0)
	cs1.Add(e, data1)
	cs2.Add(e, data2)
	cs3.Add(e, data3)
	cs4.Add(e, data4)
	cs5.Add(e, data5)
	fatEntity.ComponentMap = fatEntity.ComponentMap.Or(cs1.mask).Or(cs2.mask).Or(cs3.mask).Or(cs4.mask).Or(cs5.mask)
	// trigger events
	getComponentAddedEventsParent[T1](w).add(ctx, EntityComponentPair[T1]{
		Entity:        e,
		ComponentCopy: data1,
	})
	getComponentAddedEventsParent[T2](w).add(ctx, EntityComponentPair[T2]{
		Entity:        e,
		ComponentCopy: data2,
	})
	getComponentAddedEventsParent[T3](w).add(ctx, EntityComponentPair[T3]{
		Entity:        e,
		ComponentCopy: data3,
	})
	getComponentAddedEventsParent[T4](w).add(ctx, EntityComponentPair[T4]{
		Entity:        e,
		ComponentCopy: data4,
	})
	getComponentAddedEventsParent[T5](w).add(ctx, EntityComponentPair[T5]{
		Entity:        e,
		ComponentCopy: data5,
	})
}

func addComponent6[T1, T2, T3, T4, T5, T6 Component](ctx *Context, e Entity, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5, data6 T6) {
	w := ctx.world
	fatEntity := w.getFatEntity(e)
	assert(fatEntity != nil, "entity not found")
	cs1 := getOrCreateComponentStorage[T1](w, 0)
	cs2 := getOrCreateComponentStorage[T2](w, 0)
	cs3 := getOrCreateComponentStorage[T3](w, 0)
	cs4 := getOrCreateComponentStorage[T4](w, 0)
	cs5 := getOrCreateComponentStorage[T5](w, 0)
	cs6 := getOrCreateComponentStorage[T6](w, 0)
	cs1.Add(e, data1)
	cs2.Add(e, data2)
	cs3.Add(e, data3)
	cs4.Add(e, data4)
	cs5.Add(e, data5)
	cs6.Add(e, data6)
	fatEntity.ComponentMap = fatEntity.ComponentMap.Or(cs1.mask).Or(cs2.mask).Or(cs3.mask).Or(cs4.mask).Or(cs5.mask).Or(cs6.mask)
	// trigger events
	getComponentAddedEventsParent[T1](w).add(ctx, EntityComponentPair[T1]{
		Entity:        e,
		ComponentCopy: data1,
	})
	getComponentAddedEventsParent[T2](w).add(ctx, EntityComponentPair[T2]{
		Entity:        e,
		ComponentCopy: data2,
	})
	getComponentAddedEventsParent[T3](w).add(ctx, EntityComponentPair[T3]{
		Entity:        e,
		ComponentCopy: data3,
	})
	getComponentAddedEventsParent[T4](w).add(ctx, EntityComponentPair[T4]{
		Entity:        e,
		ComponentCopy: data4,
	})
	getComponentAddedEventsParent[T5](w).add(ctx, EntityComponentPair[T5]{
		Entity:        e,
		ComponentCopy: data5,
	})
	getComponentAddedEventsParent[T6](w).add(ctx, EntityComponentPair[T6]{
		Entity:        e,
		ComponentCopy: data6,
	})
}

func addComponent7[T1, T2, T3, T4, T5, T6, T7 Component](ctx *Context, e Entity, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5, data6 T6, data7 T7) {
	w := ctx.world
	fatEntity := w.getFatEntity(e)
	assert(fatEntity != nil, "entity not found")
	cs1 := getOrCreateComponentStorage[T1](w, 0)
	cs2 := getOrCreateComponentStorage[T2](w, 0)
	cs3 := getOrCreateComponentStorage[T3](w, 0)
	cs4 := getOrCreateComponentStorage[T4](w, 0)
	cs5 := getOrCreateComponentStorage[T5](w, 0)
	cs6 := getOrCreateComponentStorage[T6](w, 0)
	cs7 := getOrCreateComponentStorage[T7](w, 0)
	cs1.Add(e, data1)
	cs2.Add(e, data2)
	cs3.Add(e, data3)
	cs4.Add(e, data4)
	cs5.Add(e, data5)
	cs6.Add(e, data6)
	cs7.Add(e, data7)
	fatEntity.ComponentMap = fatEntity.ComponentMap.Or(cs1.mask).Or(cs2.mask).Or(cs3.mask).Or(cs4.mask).Or(cs5.mask).Or(cs6.mask).Or(cs7.mask)
	// trigger events
	getComponentAddedEventsParent[T1](w).add(ctx, EntityComponentPair[T1]{
		Entity:        e,
		ComponentCopy: data1,
	})
	getComponentAddedEventsParent[T2](w).add(ctx, EntityComponentPair[T2]{
		Entity:        e,
		ComponentCopy: data2,
	})
	getComponentAddedEventsParent[T3](w).add(ctx, EntityComponentPair[T3]{
		Entity:        e,
		ComponentCopy: data3,
	})
	getComponentAddedEventsParent[T4](w).add(ctx, EntityComponentPair[T4]{
		Entity:        e,
		ComponentCopy: data4,
	})
	getComponentAddedEventsParent[T5](w).add(ctx, EntityComponentPair[T5]{
		Entity:        e,
		ComponentCopy: data5,
	})
	getComponentAddedEventsParent[T6](w).add(ctx, EntityComponentPair[T6]{
		Entity:        e,
		ComponentCopy: data6,
	})
	getComponentAddedEventsParent[T7](w).add(ctx, EntityComponentPair[T7]{
		Entity:        e,
		ComponentCopy: data7,
	})
}

func addComponent8[T1, T2, T3, T4, T5, T6, T7, T8 Component](ctx *Context, e Entity, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5, data6 T6, data7 T7, data8 T8) {
	w := ctx.world
	fatEntity := w.getFatEntity(e)
	assert(fatEntity != nil, "entity not found")
	cs1 := getOrCreateComponentStorage[T1](w, 0)
	cs2 := getOrCreateComponentStorage[T2](w, 0)
	cs3 := getOrCreateComponentStorage[T3](w, 0)
	cs4 := getOrCreateComponentStorage[T4](w, 0)
	cs5 := getOrCreateComponentStorage[T5](w, 0)
	cs6 := getOrCreateComponentStorage[T6](w, 0)
	cs7 := getOrCreateComponentStorage[T7](w, 0)
	cs8 := getOrCreateComponentStorage[T8](w, 0)
	cs1.Add(e, data1)
	cs2.Add(e, data2)
	cs3.Add(e, data3)
	cs4.Add(e, data4)
	cs5.Add(e, data5)
	cs6.Add(e, data6)
	cs7.Add(e, data7)
	cs8.Add(e, data8)
	fatEntity.ComponentMap = fatEntity.ComponentMap.Or(cs1.mask).Or(cs2.mask).Or(cs3.mask).Or(cs4.mask).Or(cs5.mask).Or(cs6.mask).Or(cs7.mask).Or(cs8.mask)
	// trigger events
	getComponentAddedEventsParent[T1](w).add(ctx, EntityComponentPair[T1]{
		Entity:        e,
		ComponentCopy: data1,
	})
	getComponentAddedEventsParent[T2](w).add(ctx, EntityComponentPair[T2]{
		Entity:        e,
		ComponentCopy: data2,
	})
	getComponentAddedEventsParent[T3](w).add(ctx, EntityComponentPair[T3]{
		Entity:        e,
		ComponentCopy: data3,
	})
	getComponentAddedEventsParent[T4](w).add(ctx, EntityComponentPair[T4]{
		Entity:        e,
		ComponentCopy: data4,
	})
	getComponentAddedEventsParent[T5](w).add(ctx, EntityComponentPair[T5]{
		Entity:        e,
		ComponentCopy: data5,
	})
	getComponentAddedEventsParent[T6](w).add(ctx, EntityComponentPair[T6]{
		Entity:        e,
		ComponentCopy: data6,
	})
	getComponentAddedEventsParent[T7](w).add(ctx, EntityComponentPair[T7]{
		Entity:        e,
		ComponentCopy: data7,
	})
	getComponentAddedEventsParent[T8](w).add(ctx, EntityComponentPair[T8]{
		Entity:        e,
		ComponentCopy: data8,
	})
}

func addComponent9[T1, T2, T3, T4, T5, T6, T7, T8, T9 Component](ctx *Context, e Entity, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5, data6 T6, data7 T7, data8 T8, data9 T9) {
	w := ctx.world
	fatEntity := w.getFatEntity(e)
	assert(fatEntity != nil, "entity not found")
	cs1 := getOrCreateComponentStorage[T1](w, 0)
	cs2 := getOrCreateComponentStorage[T2](w, 0)
	cs3 := getOrCreateComponentStorage[T3](w, 0)
	cs4 := getOrCreateComponentStorage[T4](w, 0)
	cs5 := getOrCreateComponentStorage[T5](w, 0)
	cs6 := getOrCreateComponentStorage[T6](w, 0)
	cs7 := getOrCreateComponentStorage[T7](w, 0)
	cs8 := getOrCreateComponentStorage[T8](w, 0)
	cs9 := getOrCreateComponentStorage[T9](w, 0)
	cs1.Add(e, data1)
	cs2.Add(e, data2)
	cs3.Add(e, data3)
	cs4.Add(e, data4)
	cs5.Add(e, data5)
	cs6.Add(e, data6)
	cs7.Add(e, data7)
	cs8.Add(e, data8)
	cs9.Add(e, data9)
	fatEntity.ComponentMap = fatEntity.ComponentMap.Or(cs1.mask).Or(cs2.mask).Or(cs3.mask).Or(cs4.mask).Or(cs5.mask).Or(cs6.mask).Or(cs7.mask).Or(cs8.mask).Or(cs9.mask)
	// trigger events
	getComponentAddedEventsParent[T1](w).add(ctx, EntityComponentPair[T1]{
		Entity:        e,
		ComponentCopy: data1,
	})
	getComponentAddedEventsParent[T2](w).add(ctx, EntityComponentPair[T2]{
		Entity:        e,
		ComponentCopy: data2,
	})
	getComponentAddedEventsParent[T3](w).add(ctx, EntityComponentPair[T3]{
		Entity:        e,
		ComponentCopy: data3,
	})
	getComponentAddedEventsParent[T4](w).add(ctx, EntityComponentPair[T4]{
		Entity:        e,
		ComponentCopy: data4,
	})
	getComponentAddedEventsParent[T5](w).add(ctx, EntityComponentPair[T5]{
		Entity:        e,
		ComponentCopy: data5,
	})
	getComponentAddedEventsParent[T6](w).add(ctx, EntityComponentPair[T6]{
		Entity:        e,
		ComponentCopy: data6,
	})
	getComponentAddedEventsParent[T7](w).add(ctx, EntityComponentPair[T7]{
		Entity:        e,
		ComponentCopy: data7,
	})
	getComponentAddedEventsParent[T8](w).add(ctx, EntityComponentPair[T8]{
		Entity:        e,
		ComponentCopy: data8,
	})
	getComponentAddedEventsParent[T9](w).add(ctx, EntityComponentPair[T9]{
		Entity:        e,
		ComponentCopy: data9,
	})
}

func addComponent10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 Component](ctx *Context, e Entity, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5, data6 T6, data7 T7, data8 T8, data9 T9, data10 T10) {
	w := ctx.world
	fatEntity := w.getFatEntity(e)
	assert(fatEntity != nil, "entity not found")
	cs1 := getOrCreateComponentStorage[T1](w, 0)
	cs2 := getOrCreateComponentStorage[T2](w, 0)
	cs3 := getOrCreateComponentStorage[T3](w, 0)
	cs4 := getOrCreateComponentStorage[T4](w, 0)
	cs5 := getOrCreateComponentStorage[T5](w, 0)
	cs6 := getOrCreateComponentStorage[T6](w, 0)
	cs7 := getOrCreateComponentStorage[T7](w, 0)
	cs8 := getOrCreateComponentStorage[T8](w, 0)
	cs9 := getOrCreateComponentStorage[T9](w, 0)
	cs10 := getOrCreateComponentStorage[T10](w, 0)
	cs1.Add(e, data1)
	cs2.Add(e, data2)
	cs3.Add(e, data3)
	cs4.Add(e, data4)
	cs5.Add(e, data5)
	cs6.Add(e, data6)
	cs7.Add(e, data7)
	cs8.Add(e, data8)
	cs9.Add(e, data9)
	cs10.Add(e, data10)
	fatEntity.ComponentMap = fatEntity.ComponentMap.Or(cs1.mask).Or(cs2.mask).Or(cs3.mask).Or(cs4.mask).Or(cs5.mask).Or(cs6.mask).Or(cs7.mask).Or(cs8.mask).Or(cs9.mask).Or(cs10.mask)
	// trigger events
	getComponentAddedEventsParent[T1](w).add(ctx, EntityComponentPair[T1]{
		Entity:        e,
		ComponentCopy: data1,
	})
	getComponentAddedEventsParent[T2](w).add(ctx, EntityComponentPair[T2]{
		Entity:        e,
		ComponentCopy: data2,
	})
	getComponentAddedEventsParent[T3](w).add(ctx, EntityComponentPair[T3]{
		Entity:        e,
		ComponentCopy: data3,
	})
	getComponentAddedEventsParent[T4](w).add(ctx, EntityComponentPair[T4]{
		Entity:        e,
		ComponentCopy: data4,
	})
	getComponentAddedEventsParent[T5](w).add(ctx, EntityComponentPair[T5]{
		Entity:        e,
		ComponentCopy: data5,
	})
	getComponentAddedEventsParent[T6](w).add(ctx, EntityComponentPair[T6]{
		Entity:        e,
		ComponentCopy: data6,
	})
	getComponentAddedEventsParent[T7](w).add(ctx, EntityComponentPair[T7]{
		Entity:        e,
		ComponentCopy: data7,
	})
	getComponentAddedEventsParent[T8](w).add(ctx, EntityComponentPair[T8]{
		Entity:        e,
		ComponentCopy: data8,
	})
	getComponentAddedEventsParent[T9](w).add(ctx, EntityComponentPair[T9]{
		Entity:        e,
		ComponentCopy: data9,
	})
	getComponentAddedEventsParent[T10](w).add(ctx, EntityComponentPair[T10]{
		Entity:        e,
		ComponentCopy: data10,
	})
}
