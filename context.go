package ecs

import "reflect"

type Context struct {
	world              World
	commands           []Command
	currentSystem      *worldSystem
	isStartupSystem    bool
	currentSystemIndex int
}

func (c *Context) World() World {
	return c.world
}

func LocalResource[T any](c *Context) *T {
	if c.isStartupSystem {
		panic("LocalResource is not allowed in startup systems")
	}
	var zv T
	tm := typeMapKeyOf(reflect.TypeOf(zv))
	x := c.currentSystem.LocalResources[tm]
	if x == nil {
		zvp := &zv
		if vi, ok := any(zvp).(WorldIniter); ok {
			vi.Init(c.world)
		}
		c.currentSystem.LocalResources[tm] = zvp
	}
	return (c.currentSystem.LocalResources[tm].(*T))
}

type EntityCommandCallback func(ctx *Context, e Entity)

func Spawn[T Component](c *Context, data T, actions ...EntityCommandCallback) {
	c.commands = append(c.commands, newSpawnCommand(c.world, data))
}

func Spawn2[T1, T2 Component](c *Context, data1 T1, data2 T2) {
	c.commands = append(c.commands, newSpawn2Command(c.world, data1, data2))
}

func Spawn3[T1, T2, T3 Component](c *Context, data1 T1, data2 T2, data3 T3) {
	c.commands = append(c.commands, newSpawn3Command(c.world, data1, data2, data3))
}

func Spawn4[T1, T2, T3, T4 Component](c *Context, data1 T1, data2 T2, data3 T3, data4 T4) {
	c.commands = append(c.commands, newSpawn4Command(c.world, data1, data2, data3, data4))
}

func Spawn5[T1, T2, T3, T4, T5 Component](c *Context, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5) {
	c.commands = append(c.commands, newSpawn5Command(c.world, data1, data2, data3, data4, data5))
}

func Spawn6[T1, T2, T3, T4, T5, T6 Component](c *Context, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5, data6 T6) {
	c.commands = append(c.commands, newSpawn6Command(c.world, data1, data2, data3, data4, data5, data6))
}

func Spawn7[T1, T2, T3, T4, T5, T6, T7 Component](c *Context, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5, data6 T6, data7 T7) {
	c.commands = append(c.commands, newSpawn7Command(c.world, data1, data2, data3, data4, data5, data6, data7))
}

func Spawn8[T1, T2, T3, T4, T5, T6, T7, T8 Component](c *Context, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5, data6 T6, data7 T7, data8 T8) {
	c.commands = append(c.commands, newSpawn8Command(c.world, data1, data2, data3, data4, data5, data6, data7, data8))
}

func Spawn9[T1, T2, T3, T4, T5, T6, T7, T8, T9 Component](c *Context, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5, data6 T6, data7 T7, data8 T8, data9 T9) {
	c.commands = append(c.commands, newSpawn9Command(c.world, data1, data2, data3, data4, data5, data6, data7, data8, data9))
}

func Spawn10[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 Component](c *Context, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5, data6 T6, data7 T7, data8 T8, data9 T9, data10 T10) {
	c.commands = append(c.commands, newSpawn10Command(c.world, data1, data2, data3, data4, data5, data6, data7, data8, data9, data10))
}

func RemoveEntity(ctx *Context, e Entity) {
	ctx.commands = append(ctx.commands, func(ctx *Context) {
		ctx.world.removeEntity(ctx, e)
		//the component removed event is called inside the removeEntity function
	})
}

func RemoveComponent[T Component](ctx *Context, e Entity) {
	ctx.commands = append(ctx.commands, func(ctx *Context) {
		removeComponent[T](ctx, e)
	})
}

func AddComponent[T Component](ctx *Context, e Entity, data T, actions ...EntityCommandCallback) {
	ctx.commands = append(ctx.commands, func(ctx *Context) {
		addComponent(ctx, e, data)
		execSpawnCallbacks(ctx, e, actions...)
	})
}

func (ctx *Context) run() {
	for _, cmd := range ctx.commands {
		cmd(ctx)
	}
	//TODO: reorganize entities (if needed) after all commands are executed
}

type Command func(parent *Context)

func execSpawnCallbacks(parentctx *Context, e Entity, actions ...EntityCommandCallback) {
	if len(actions) < 1 {
		return
	}
	ctxchild := &Context{
		world:              parentctx.world,
		commands:           make([]Command, 0),
		currentSystem:      parentctx.currentSystem,
		isStartupSystem:    parentctx.isStartupSystem,
		currentSystemIndex: parentctx.currentSystemIndex,
	}
	for _, action := range actions {
		action(ctxchild, e)
	}
	ctxchild.run()
	ctxchild.commands = nil
}

func newSpawnCommand[T Component](w World, data T, actions ...EntityCommandCallback) Command {
	return func(parent *Context) {
		e := w.newEntity()
		addComponent(parent, e, data)
		execSpawnCallbacks(parent, e, actions...)
	}
}

func newSpawn2Command[T1, T2 Component](w World, data1 T1, data2 T2, actions ...EntityCommandCallback) Command {
	return func(parent *Context) {
		e := w.newEntity()
		addComponent2(parent, e, data1, data2)
		execSpawnCallbacks(parent, e, actions...)
	}
}

func newSpawn3Command[T1, T2, T3 Component](w World, data1 T1, data2 T2, data3 T3, actions ...EntityCommandCallback) Command {
	return func(parent *Context) {
		e := w.newEntity()
		addComponent3(parent, e, data1, data2, data3)
		execSpawnCallbacks(parent, e, actions...)
	}
}

func newSpawn4Command[T1, T2, T3, T4 Component](w World, data1 T1, data2 T2, data3 T3, data4 T4, actions ...EntityCommandCallback) Command {
	return func(parent *Context) {
		e := w.newEntity()
		addComponent4(parent, e, data1, data2, data3, data4)
		execSpawnCallbacks(parent, e, actions...)
	}
}

func newSpawn5Command[T1, T2, T3, T4, T5 Component](w World, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5, actions ...EntityCommandCallback) Command {
	return func(parent *Context) {
		e := w.newEntity()
		addComponent5(parent, e, data1, data2, data3, data4, data5)
		execSpawnCallbacks(parent, e, actions...)
	}
}

func newSpawn6Command[T1, T2, T3, T4, T5, T6 Component](w World, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5, data6 T6, actions ...EntityCommandCallback) Command {
	return func(parent *Context) {
		e := w.newEntity()
		addComponent6(parent, e, data1, data2, data3, data4, data5, data6)
		execSpawnCallbacks(parent, e, actions...)
	}
}

func newSpawn7Command[T1, T2, T3, T4, T5, T6, T7 Component](w World, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5, data6 T6, data7 T7, actions ...EntityCommandCallback) Command {
	return func(parent *Context) {
		e := w.newEntity()
		addComponent7(parent, e, data1, data2, data3, data4, data5, data6, data7)
		execSpawnCallbacks(parent, e, actions...)
	}
}

func newSpawn8Command[T1, T2, T3, T4, T5, T6, T7, T8 Component](w World, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5, data6 T6, data7 T7, data8 T8, actions ...EntityCommandCallback) Command {
	return func(parent *Context) {
		e := w.newEntity()
		addComponent8(parent, e, data1, data2, data3, data4, data5, data6, data7, data8)
		execSpawnCallbacks(parent, e, actions...)
	}
}

func newSpawn9Command[T1, T2, T3, T4, T5, T6, T7, T8, T9 Component](w World, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5, data6 T6, data7 T7, data8 T8, data9 T9, actions ...EntityCommandCallback) Command {
	return func(parent *Context) {
		e := w.newEntity()
		addComponent9(parent, e, data1, data2, data3, data4, data5, data6, data7, data8, data9)
		execSpawnCallbacks(parent, e, actions...)
	}
}

func newSpawn10Command[T1, T2, T3, T4, T5, T6, T7, T8, T9, T10 Component](w World, data1 T1, data2 T2, data3 T3, data4 T4, data5 T5, data6 T6, data7 T7, data8 T8, data9 T9, data10 T10, actions ...EntityCommandCallback) Command {
	return func(parent *Context) {
		e := w.newEntity()
		addComponent10(parent, e, data1, data2, data3, data4, data5, data6, data7, data8, data9, data10)
		execSpawnCallbacks(parent, e, actions...)
	}
}
