package simple

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

//
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

// Position component
//
// ecs:component
// uuid:3DF7F486-807D-4CE8-A187-37CED338137B
type Position struct {
	X float64
	Y float64
}

// Rotation component
//
// ecs:component
// uuid:55AFAC07-F446-4DBE-B963-413B0C38E72B
type Rotation struct {
	Radians float64
}

// Velocity component
//
// ecs:component
// the uuid for this component is generated automatically
type Velocity struct {
	X       float64
	Y       float64
	Radians float64
}

//go:generate go run ../../cmd/ecsgen/main.go

// Update MovementSystem matches
//
// ecs:system
// uuid:E7D2FB64-2E98-4A14-8DE6-F088DE2AC3FB
// name:MovementSystem
// components: Position, Rotation, Velocity
//
// member: lastdt float64
//
// entityadded: onEntityAdded
// entityremoved: onEntityRemoved
// componentwillresize: onComponentWillResize
// componentresized: onComponentResized
// setup: onSetup
func (s *MovementSystem) Update(dt float64) {
	for _, v := range s.V().Matches() {
		v.Position.X += v.Velocity.X
		v.Position.Y += v.Velocity.Y
		v.Rotation.Radians += v.Velocity.Radians
	}
}

func (s *MovementSystem) onEntityAdded(e ecs.Entity) {
	println("onEntityAdded")
}

func (s *MovementSystem) onEntityRemoved(e ecs.Entity) {
	println("onEntityRemoved")
}

func (s *MovementSystem) onComponentWillResize(cflag ecs.Flag) {
	println("onComponentWillResize")
}

func (s *MovementSystem) onComponentResized(cflag ecs.Flag) {
	println("onComponentResized")
}

func (s *MovementSystem) onSetup(w ecs.World) {
	println("onSetup")
}
