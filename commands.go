package ecs

type Commands struct {
	world World
	list  []Command
}

func (c *Commands) World() World {
	return c.world
}

func Spawn[T Component](c *Commands, data T) {
	c.list = append(c.list, newSpawnCommand(c.world, data))
}

func Spawn2[T1, T2 Component](c *Commands, data1 T1, data2 T2) {
	c.list = append(c.list, newSpawn2Command(c.world, data1, data2))
}

func Spawn3[T1, T2, T3 Component](c *Commands, data1 T1, data2 T2, data3 T3) {
	c.list = append(c.list, newSpawn3Command(c.world, data1, data2, data3))
}

func Spawn4[T1, T2, T3, T4 Component](c *Commands, data1 T1, data2 T2, data3 T3, data4 T4) {
	c.list = append(c.list, newSpawn4Command(c.world, data1, data2, data3, data4))
}

func Spawn5[T1, T2, T3, T4, T5 Component](c *Commands, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5) {
	c.list = append(c.list, newSpawn5Command(c.world, data1, data2, data3, data4, data5))
}

func Spawn6[T1, T2, T3, T4, T5, T6 Component](c *Commands, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5, data6 T6) {
	c.list = append(c.list, newSpawn6Command(c.world, data1, data2, data3, data4, data5, data6))
}

func Spawn7[T1, T2, T3, T4, T5, T6, T7 Component](c *Commands, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5, data6 T6, data7 T7) {
	c.list = append(c.list, newSpawn7Command(c.world, data1, data2, data3, data4, data5, data6, data7))
}

func Spawn8[T1, T2, T3, T4, T5, T6, T7, T8 Component](c *Commands, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5, data6 T6, data7 T7, data8 T8) {
	c.list = append(c.list, newSpawn8Command(c.world, data1, data2, data3, data4, data5, data6, data7, data8))
}

func Spawn9[T1, T2, T3, T4, T5, T6, T7, T8, T9 Component](c *Commands, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5, data6 T6, data7 T7, data8 T8, data9 T9) {
	c.list = append(c.list, newSpawn9Command(c.world, data1, data2, data3, data4, data5, data6, data7, data8, data9))
}

func Spawn10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 Component](c *Commands, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5, data6 T6, data7 T7, data8 T8, data9 T9, data10 T10) {
	c.list = append(c.list, newSpawn10Command(c.world, data1, data2, data3, data4, data5, data6, data7, data8, data9, data10))
}

func RemoveEntity(c *Commands, e Entity) {
	c.list = append(c.list, func() {
		c.world.removeEntity(e)
	})
}

func RemoveComponent[T Component](c *Commands, e Entity) {
	c.list = append(c.list, func() {
		removeComponent[T](c.world, e)
	})
}

func AddComponent[T Component](c *Commands, e Entity, data T) {
	c.list = append(c.list, func() {
		addComponent(c.world, e, data)
	})
}

func (c *Commands) run() {
	for _, cmd := range c.list {
		cmd()
	}
	//TODO: reorganize entities (if needed) after all commands are executed
}

type Command func()

func newSpawnCommand[T Component](w World, data T) Command {
	return func() {
		e := w.newEntity()
		addComponent(w, e, data)
	}
}

func newSpawn2Command[T1, T2 Component](w World, data1 T1, data2 T2) Command {
	return func() {
		e := w.newEntity()
		addComponent2(w, e, data1, data2)
	}
}

func newSpawn3Command[T1, T2, T3 Component](w World, data1 T1, data2 T2, data3 T3) Command {
	return func() {
		e := w.newEntity()
		addComponent3(w, e, data1, data2, data3)
	}
}

func newSpawn4Command[T1, T2, T3, T4 Component](w World, data1 T1, data2 T2, data3 T3, data4 T4) Command {
	return func() {
		e := w.newEntity()
		addComponent4(w, e, data1, data2, data3, data4)
	}
}

func newSpawn5Command[T1, T2, T3, T4, T5 Component](w World, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5) Command {
	return func() {
		e := w.newEntity()
		addComponent5(w, e, data1, data2, data3, data4, data5)
	}
}

func newSpawn6Command[T1, T2, T3, T4, T5, T6 Component](w World, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5, data6 T6) Command {
	return func() {
		e := w.newEntity()
		addComponent6(w, e, data1, data2, data3, data4, data5, data6)
	}
}

func newSpawn7Command[T1, T2, T3, T4, T5, T6, T7 Component](w World, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5, data6 T6, data7 T7) Command {
	return func() {
		e := w.newEntity()
		addComponent7(w, e, data1, data2, data3, data4, data5, data6, data7)
	}
}

func newSpawn8Command[T1, T2, T3, T4, T5, T6, T7, T8 Component](w World, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5, data6 T6, data7 T7, data8 T8) Command {
	return func() {
		e := w.newEntity()
		addComponent8(w, e, data1, data2, data3, data4, data5, data6, data7, data8)
	}
}

func newSpawn9Command[T1, T2, T3, T4, T5, T6, T7, T8, T9 Component](w World, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5, data6 T6, data7 T7, data8 T8, data9 T9) Command {
	return func() {
		e := w.newEntity()
		addComponent9(w, e, data1, data2, data3, data4, data5, data6, data7, data8, data9)
	}
}

func newSpawn10Command[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 Component](w World, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5, data6 T6, data7 T7, data8 T8, data9 T9, data10 T10) Command {
	return func() {
		e := w.newEntity()
		addComponent10(w, e, data1, data2, data3, data4, data5, data6, data7, data8, data9, data10)
	}
}
