# Entity Component System

A fast, code generate ECS (no more interface{}). Game Engine agnostic.

[![GoDoc](https://godoc.org/github.com/gabstv/ecs?status.svg)](https://godoc.org/github.com/gabstv/ecs)

`go get github.com/gabstv/ecs/v2/cmd/ecsgen`

Example:

```go
package mycomponents

import "github.com/gabstv/ecs/v2"


type World struct {
	ecs.World
}

func NewWorld() *World {
	base := ecs.NewWorld()
	ecs.RegisterWorldDefaults(base)
	return &World{
		World: base,
	}
}

// Update all systems (sorted by priority)
func (w *World) Update(dt float64) {
	w.EachSystem(func(s ecs.System) bool {
		s.(System).Update(dt)
		return true
	})
}

type System interface {
	ecs.System
	Update(dt float64)
}

// ecsgen will create the component+system logic:

//go:generate go run ecsgen

// Position component data
//
// ecs:component
// uuid:9F414CAD-4C1B-49B2-980E-0A61302AD5DE
type Position struct {
	X float64
	Y float64
}

// Velocity component
//
// ecs:component
// the uuid for this component is generated automatically
type Velocity struct {
	X       float64
	Y       float64
}

// Update MovementSystem matches
//
// ecs:system
// uuid:43838027-AA12-4AD2-9F09-5DCBDA589779
// name:MovementSystem
// components: Position, Velocity
func (s *MovementSystem) Update(dt float64) {
	for _, v := range s.V().Matches() {
		v.Position.X += v.Velocity.X
		v.Position.Y += v.Velocity.Y
	}
}

```