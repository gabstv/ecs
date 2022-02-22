package main

import (
	"fmt"

	"github.com/gabstv/ecs/v3"
)

// system flags to filten which systems to run at a time
const (
	SysFlagUpdate  = 1
	SysFlagExample = 2
)

type Position struct {
	X, Y float64
}

// Pkg is required for the component registry to avoid using reflection
func (_ Position) Pkg() string {
	return "main.Position"
}

type Speed struct {
	X, Y float64
}

// Pkg is required for the component registry to avoid using reflection
func (_ Speed) Pkg() string {
	return "main.Speed"
}

type Accel struct {
	X, Y float64
}

// Pkg is required for the component registry to avoid using reflection
func (_ Accel) Pkg() string {
	return "main.Accel"
}

func main() {
	world := ecs.NewWorld()
	setupSystems(world)

	e1 := world.NewEntity()
	e2 := world.NewEntity()

	ecs.Set(world, e1, Position{
		X: 10,
		Y: 10,
	})
	ecs.Set(world, e1, Speed{
		X: 1,
		Y: .5,
	})
	ecs.Set(world, e1, Accel{
		X: .1,
		Y: .05,
	})
	ecs.Set(world, e2, Position{
		X: 100,
		Y: -50,
	})
	ecs.Set(world, e2, Speed{
		X: -1,
		Y: 1,
	})
	// no accel for e2 (and that's ok)

	world.Step()                // run ALL systems once
	world.StepF(SysFlagUpdate)  // run all systems with this flag
	world.StepF(SysFlagExample) // run all systems with this flag

	ecs.Apply(world, e1, func(pos *Position) {
		fmt.Println("e1 Position:", pos)
	})
	ecs.Apply(world, e2, func(pos *Position) {
		fmt.Println("e2 Position:", pos)
	})

	// remove e2's Speed component
	ecs.RemoveComponent[Speed](world, e2)

	world.Step()
	world.Step()
	world.Step()

	ecs.Apply(world, e1, func(pos *Position) {
		fmt.Println("e1 Position:", pos)
	})
	ecs.Apply(world, e2, func(pos *Position) {
		fmt.Println("e2 Position:", pos)
	})
}

func setupSystems(world *ecs.World) {
	// acceleration will be applied here
	accelSys := ecs.NewSystem2[Speed, Accel](0, world)
	accelSys.SetFlag(SysFlagUpdate)
	accelSys.Run = func(view *ecs.View2[Speed, Accel]) {
		view.Each(func(_ ecs.Entity, speed *Speed, accel *Accel) {
			speed.X += accel.X
			speed.Y += accel.Y
		})
	}

	// speed will be applied here
	speedSys := ecs.NewSystem2[Position, Speed](1, world)
	speedSys.SetFlag(SysFlagUpdate)
	speedSys.Run = func(view *ecs.View2[Position, Speed]) {
		view.Each(func(_ ecs.Entity, pos *Position, speed *Speed) {
			pos.X += speed.X
			pos.Y += speed.Y
		})
	}
}
