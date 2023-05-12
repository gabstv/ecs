package ecs

func addComponent[T Component](w World, e Entity, data T) {
	fatEntity := w.getFatEntity(e)
	assert(fatEntity != nil, "entity not found")
	cs := getOrCreateComponentStorage[T](w)
	cs.Add(e, data)
	fatEntity.ComponentMap = fatEntity.ComponentMap.Or(cs.mask)
	// trigger events
	getComponentAddedEventsParent[T](w).add(EntityComponentPair[T]{
		Entity:        e,
		ComponentCopy: data,
	})
}

func addComponent2[T1, T2 Component](w World, e Entity, data1 T1, data2 T2) {
	fatEntity := w.getFatEntity(e)
	assert(fatEntity != nil, "entity not found")
	cs1 := getOrCreateComponentStorage[T1](w)
	cs2 := getOrCreateComponentStorage[T2](w)
	cs1.Add(e, data1)
	cs2.Add(e, data2)
	fatEntity.ComponentMap = fatEntity.ComponentMap.Or(cs1.mask).Or(cs2.mask)
	// trigger events
	getComponentAddedEventsParent[T1](w).add(EntityComponentPair[T1]{
		Entity:        e,
		ComponentCopy: data1,
	})
	getComponentAddedEventsParent[T2](w).add(EntityComponentPair[T2]{
		Entity:        e,
		ComponentCopy: data2,
	})
}

func addComponent3[T1, T2, T3 Component](w World, e Entity, data1 T1, data2 T2, data3 T3) {
	fatEntity := w.getFatEntity(e)
	assert(fatEntity != nil, "entity not found")
	cs1 := getOrCreateComponentStorage[T1](w)
	cs2 := getOrCreateComponentStorage[T2](w)
	cs3 := getOrCreateComponentStorage[T3](w)
	cs1.Add(e, data1)
	cs2.Add(e, data2)
	cs3.Add(e, data3)
	fatEntity.ComponentMap = fatEntity.ComponentMap.Or(cs1.mask).Or(cs2.mask).Or(cs3.mask)
	// trigger events
	getComponentAddedEventsParent[T1](w).add(EntityComponentPair[T1]{
		Entity:        e,
		ComponentCopy: data1,
	})
	getComponentAddedEventsParent[T2](w).add(EntityComponentPair[T2]{
		Entity:        e,
		ComponentCopy: data2,
	})
	getComponentAddedEventsParent[T3](w).add(EntityComponentPair[T3]{
		Entity:        e,
		ComponentCopy: data3,
	})
}

func addComponent4[T1, T2, T3, T4 Component](w World, e Entity, data1 T1, data2 T2, data3 T3, data4 T4) {
	fatEntity := w.getFatEntity(e)
	assert(fatEntity != nil, "entity not found")
	cs1 := getOrCreateComponentStorage[T1](w)
	cs2 := getOrCreateComponentStorage[T2](w)
	cs3 := getOrCreateComponentStorage[T3](w)
	cs4 := getOrCreateComponentStorage[T4](w)
	cs1.Add(e, data1)
	cs2.Add(e, data2)
	cs3.Add(e, data3)
	cs4.Add(e, data4)
	fatEntity.ComponentMap = fatEntity.ComponentMap.Or(cs1.mask).Or(cs2.mask).Or(cs3.mask).Or(cs4.mask)
	// trigger events
	getComponentAddedEventsParent[T1](w).add(EntityComponentPair[T1]{
		Entity:        e,
		ComponentCopy: data1,
	})
	getComponentAddedEventsParent[T2](w).add(EntityComponentPair[T2]{
		Entity:        e,
		ComponentCopy: data2,
	})
	getComponentAddedEventsParent[T3](w).add(EntityComponentPair[T3]{
		Entity:        e,
		ComponentCopy: data3,
	})
	getComponentAddedEventsParent[T4](w).add(EntityComponentPair[T4]{
		Entity:        e,
		ComponentCopy: data4,
	})
}

func addComponent5[T1, T2, T3, T4, T5 Component](w World, e Entity, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5) {
	fatEntity := w.getFatEntity(e)
	assert(fatEntity != nil, "entity not found")
	cs1 := getOrCreateComponentStorage[T1](w)
	cs2 := getOrCreateComponentStorage[T2](w)
	cs3 := getOrCreateComponentStorage[T3](w)
	cs4 := getOrCreateComponentStorage[T4](w)
	cs5 := getOrCreateComponentStorage[T5](w)
	cs1.Add(e, data1)
	cs2.Add(e, data2)
	cs3.Add(e, data3)
	cs4.Add(e, data4)
	cs5.Add(e, data5)
	fatEntity.ComponentMap = fatEntity.ComponentMap.Or(cs1.mask).Or(cs2.mask).Or(cs3.mask).Or(cs4.mask).Or(cs5.mask)
	// trigger events
	getComponentAddedEventsParent[T1](w).add(EntityComponentPair[T1]{
		Entity:        e,
		ComponentCopy: data1,
	})
	getComponentAddedEventsParent[T2](w).add(EntityComponentPair[T2]{
		Entity:        e,
		ComponentCopy: data2,
	})
	getComponentAddedEventsParent[T3](w).add(EntityComponentPair[T3]{
		Entity:        e,
		ComponentCopy: data3,
	})
	getComponentAddedEventsParent[T4](w).add(EntityComponentPair[T4]{
		Entity:        e,
		ComponentCopy: data4,
	})
	getComponentAddedEventsParent[T5](w).add(EntityComponentPair[T5]{
		Entity:        e,
		ComponentCopy: data5,
	})
}

func addComponent6[T1, T2, T3, T4, T5, T6 Component](w World, e Entity, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5, data6 T6) {
	fatEntity := w.getFatEntity(e)
	assert(fatEntity != nil, "entity not found")
	cs1 := getOrCreateComponentStorage[T1](w)
	cs2 := getOrCreateComponentStorage[T2](w)
	cs3 := getOrCreateComponentStorage[T3](w)
	cs4 := getOrCreateComponentStorage[T4](w)
	cs5 := getOrCreateComponentStorage[T5](w)
	cs6 := getOrCreateComponentStorage[T6](w)
	cs1.Add(e, data1)
	cs2.Add(e, data2)
	cs3.Add(e, data3)
	cs4.Add(e, data4)
	cs5.Add(e, data5)
	cs6.Add(e, data6)
	fatEntity.ComponentMap = fatEntity.ComponentMap.Or(cs1.mask).Or(cs2.mask).Or(cs3.mask).Or(cs4.mask).Or(cs5.mask).Or(cs6.mask)
	// trigger events
	getComponentAddedEventsParent[T1](w).add(EntityComponentPair[T1]{
		Entity:        e,
		ComponentCopy: data1,
	})
	getComponentAddedEventsParent[T2](w).add(EntityComponentPair[T2]{
		Entity:        e,
		ComponentCopy: data2,
	})
	getComponentAddedEventsParent[T3](w).add(EntityComponentPair[T3]{
		Entity:        e,
		ComponentCopy: data3,
	})
	getComponentAddedEventsParent[T4](w).add(EntityComponentPair[T4]{
		Entity:        e,
		ComponentCopy: data4,
	})
	getComponentAddedEventsParent[T5](w).add(EntityComponentPair[T5]{
		Entity:        e,
		ComponentCopy: data5,
	})
	getComponentAddedEventsParent[T6](w).add(EntityComponentPair[T6]{
		Entity:        e,
		ComponentCopy: data6,
	})
}

func addComponent7[T1, T2, T3, T4, T5, T6, T7 Component](w World, e Entity, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5, data6 T6, data7 T7) {
	fatEntity := w.getFatEntity(e)
	assert(fatEntity != nil, "entity not found")
	cs1 := getOrCreateComponentStorage[T1](w)
	cs2 := getOrCreateComponentStorage[T2](w)
	cs3 := getOrCreateComponentStorage[T3](w)
	cs4 := getOrCreateComponentStorage[T4](w)
	cs5 := getOrCreateComponentStorage[T5](w)
	cs6 := getOrCreateComponentStorage[T6](w)
	cs7 := getOrCreateComponentStorage[T7](w)
	cs1.Add(e, data1)
	cs2.Add(e, data2)
	cs3.Add(e, data3)
	cs4.Add(e, data4)
	cs5.Add(e, data5)
	cs6.Add(e, data6)
	cs7.Add(e, data7)
	fatEntity.ComponentMap = fatEntity.ComponentMap.Or(cs1.mask).Or(cs2.mask).Or(cs3.mask).Or(cs4.mask).Or(cs5.mask).Or(cs6.mask).Or(cs7.mask)
	// trigger events
	getComponentAddedEventsParent[T1](w).add(EntityComponentPair[T1]{
		Entity:        e,
		ComponentCopy: data1,
	})
	getComponentAddedEventsParent[T2](w).add(EntityComponentPair[T2]{
		Entity:        e,
		ComponentCopy: data2,
	})
	getComponentAddedEventsParent[T3](w).add(EntityComponentPair[T3]{
		Entity:        e,
		ComponentCopy: data3,
	})
	getComponentAddedEventsParent[T4](w).add(EntityComponentPair[T4]{
		Entity:        e,
		ComponentCopy: data4,
	})
	getComponentAddedEventsParent[T5](w).add(EntityComponentPair[T5]{
		Entity:        e,
		ComponentCopy: data5,
	})
	getComponentAddedEventsParent[T6](w).add(EntityComponentPair[T6]{
		Entity:        e,
		ComponentCopy: data6,
	})
	getComponentAddedEventsParent[T7](w).add(EntityComponentPair[T7]{
		Entity:        e,
		ComponentCopy: data7,
	})
}

func addComponent8[T1, T2, T3, T4, T5, T6, T7, T8 Component](w World, e Entity, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5, data6 T6, data7 T7, data8 T8) {
	fatEntity := w.getFatEntity(e)
	assert(fatEntity != nil, "entity not found")
	cs1 := getOrCreateComponentStorage[T1](w)
	cs2 := getOrCreateComponentStorage[T2](w)
	cs3 := getOrCreateComponentStorage[T3](w)
	cs4 := getOrCreateComponentStorage[T4](w)
	cs5 := getOrCreateComponentStorage[T5](w)
	cs6 := getOrCreateComponentStorage[T6](w)
	cs7 := getOrCreateComponentStorage[T7](w)
	cs8 := getOrCreateComponentStorage[T8](w)
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
	getComponentAddedEventsParent[T1](w).add(EntityComponentPair[T1]{
		Entity:        e,
		ComponentCopy: data1,
	})
	getComponentAddedEventsParent[T2](w).add(EntityComponentPair[T2]{
		Entity:        e,
		ComponentCopy: data2,
	})
	getComponentAddedEventsParent[T3](w).add(EntityComponentPair[T3]{
		Entity:        e,
		ComponentCopy: data3,
	})
	getComponentAddedEventsParent[T4](w).add(EntityComponentPair[T4]{
		Entity:        e,
		ComponentCopy: data4,
	})
	getComponentAddedEventsParent[T5](w).add(EntityComponentPair[T5]{
		Entity:        e,
		ComponentCopy: data5,
	})
	getComponentAddedEventsParent[T6](w).add(EntityComponentPair[T6]{
		Entity:        e,
		ComponentCopy: data6,
	})
	getComponentAddedEventsParent[T7](w).add(EntityComponentPair[T7]{
		Entity:        e,
		ComponentCopy: data7,
	})
	getComponentAddedEventsParent[T8](w).add(EntityComponentPair[T8]{
		Entity:        e,
		ComponentCopy: data8,
	})
}

func addComponent9[T1, T2, T3, T4, T5, T6, T7, T8, T9 Component](w World, e Entity, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5, data6 T6, data7 T7, data8 T8, data9 T9) {
	fatEntity := w.getFatEntity(e)
	assert(fatEntity != nil, "entity not found")
	cs1 := getOrCreateComponentStorage[T1](w)
	cs2 := getOrCreateComponentStorage[T2](w)
	cs3 := getOrCreateComponentStorage[T3](w)
	cs4 := getOrCreateComponentStorage[T4](w)
	cs5 := getOrCreateComponentStorage[T5](w)
	cs6 := getOrCreateComponentStorage[T6](w)
	cs7 := getOrCreateComponentStorage[T7](w)
	cs8 := getOrCreateComponentStorage[T8](w)
	cs9 := getOrCreateComponentStorage[T9](w)
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
	getComponentAddedEventsParent[T1](w).add(EntityComponentPair[T1]{
		Entity:        e,
		ComponentCopy: data1,
	})
	getComponentAddedEventsParent[T2](w).add(EntityComponentPair[T2]{
		Entity:        e,
		ComponentCopy: data2,
	})
	getComponentAddedEventsParent[T3](w).add(EntityComponentPair[T3]{
		Entity:        e,
		ComponentCopy: data3,
	})
	getComponentAddedEventsParent[T4](w).add(EntityComponentPair[T4]{
		Entity:        e,
		ComponentCopy: data4,
	})
	getComponentAddedEventsParent[T5](w).add(EntityComponentPair[T5]{
		Entity:        e,
		ComponentCopy: data5,
	})
	getComponentAddedEventsParent[T6](w).add(EntityComponentPair[T6]{
		Entity:        e,
		ComponentCopy: data6,
	})
	getComponentAddedEventsParent[T7](w).add(EntityComponentPair[T7]{
		Entity:        e,
		ComponentCopy: data7,
	})
	getComponentAddedEventsParent[T8](w).add(EntityComponentPair[T8]{
		Entity:        e,
		ComponentCopy: data8,
	})
	getComponentAddedEventsParent[T9](w).add(EntityComponentPair[T9]{
		Entity:        e,
		ComponentCopy: data9,
	})
}

func addComponent10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 Component](w World, e Entity, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5, data6 T6, data7 T7, data8 T8, data9 T9, data10 T10) {
	fatEntity := w.getFatEntity(e)
	assert(fatEntity != nil, "entity not found")
	cs1 := getOrCreateComponentStorage[T1](w)
	cs2 := getOrCreateComponentStorage[T2](w)
	cs3 := getOrCreateComponentStorage[T3](w)
	cs4 := getOrCreateComponentStorage[T4](w)
	cs5 := getOrCreateComponentStorage[T5](w)
	cs6 := getOrCreateComponentStorage[T6](w)
	cs7 := getOrCreateComponentStorage[T7](w)
	cs8 := getOrCreateComponentStorage[T8](w)
	cs9 := getOrCreateComponentStorage[T9](w)
	cs10 := getOrCreateComponentStorage[T10](w)
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
	getComponentAddedEventsParent[T1](w).add(EntityComponentPair[T1]{
		Entity:        e,
		ComponentCopy: data1,
	})
	getComponentAddedEventsParent[T2](w).add(EntityComponentPair[T2]{
		Entity:        e,
		ComponentCopy: data2,
	})
	getComponentAddedEventsParent[T3](w).add(EntityComponentPair[T3]{
		Entity:        e,
		ComponentCopy: data3,
	})
	getComponentAddedEventsParent[T4](w).add(EntityComponentPair[T4]{
		Entity:        e,
		ComponentCopy: data4,
	})
	getComponentAddedEventsParent[T5](w).add(EntityComponentPair[T5]{
		Entity:        e,
		ComponentCopy: data5,
	})
	getComponentAddedEventsParent[T6](w).add(EntityComponentPair[T6]{
		Entity:        e,
		ComponentCopy: data6,
	})
	getComponentAddedEventsParent[T7](w).add(EntityComponentPair[T7]{
		Entity:        e,
		ComponentCopy: data7,
	})
	getComponentAddedEventsParent[T8](w).add(EntityComponentPair[T8]{
		Entity:        e,
		ComponentCopy: data8,
	})
	getComponentAddedEventsParent[T9](w).add(EntityComponentPair[T9]{
		Entity:        e,
		ComponentCopy: data9,
	})
	getComponentAddedEventsParent[T10](w).add(EntityComponentPair[T10]{
		Entity:        e,
		ComponentCopy: data10,
	})
}
