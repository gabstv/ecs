# Entity Component System

A fast, zero reflection ECS (no more interface{}). Game Engine agnostic.

[![GoDoc](https://godoc.org/github.com/gabstv/ecs?status.svg)](https://godoc.org/github.com/gabstv/ecs)

`go get github.com/gabstv/ecs/v3`

Example:

```go
package mygame

import "github.com/gabstv/ecs/v3"

type Position struct {
	X, Y float64
}

// Pkg is required for the component registry to avoid using reflection
func (_ Position) Pkg() string {
	// ecs uses this string to identify which component registry to use
	return "main.Position" 
}

type Speed struct {
	X, Y float64
}

// Pkg is required for the component registry to avoid using reflection
func (_ Speed) Pkg() string {
	return "main.Speed"
}

func main() {
	world := ecs.NewWorld()

	// speed will be applied here
	exampleSys := ecs.NewSystem2[Position, Speed](1, world)
	exampleSys.Run = func(view *ecs.View2[Position, Speed]) {
		// you can get a reference to global (costly) variables
		// before running the view loop, like obtaining the delta time
		view.Each(func(_ ecs.Entity, pos *Position, speed *Speed) {
			pos.X += speed.X
			pos.Y += speed.Y
		})
	}

	e1 := world.NewEntity()

	ecs.Set(world, e1, Position{
		X: 10,
		Y: 20,
	})
	ecs.Set(world, e1, Speed{
		X: 1,
		Y: .5,
	})

	world.Step() // run all systems once
}

```

For a more detailed example, check the `example` folder.