package ecs

func addComponent[T Component](w World, e Entity, data T) {
	fatEntity := w.getFatEntity(e)
	assert(fatEntity != nil, "entity not found")
	cs := getOrCreateComponentStorage[T](w)
	cs.Add(e, data)
	fatEntity.ComponentMap = fatEntity.ComponentMap.Or(cs.mask)
}

func addComponent2[T1, T2 Component](w World, e Entity, data1 T1, data2 T2) {
	fatEntity := w.getFatEntity(e)
	assert(fatEntity != nil, "entity not found")
	cs1 := getOrCreateComponentStorage[T1](w)
	cs2 := getOrCreateComponentStorage[T2](w)
	cs1.Add(e, data1)
	cs2.Add(e, data2)
	fatEntity.ComponentMap = fatEntity.ComponentMap.Or(cs1.mask).Or(cs2.mask)
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
}