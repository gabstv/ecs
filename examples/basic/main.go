package main

import (
	"fmt"

	"github.com/gabstv/ecs/v4"
)

func main() {
	world := ecs.NewWorld()
	ecs.AddSystem(world, moveAll)
	ecs.AddStartupSystem(world, setup)
	world.Step()
	world.Step()
	q := ecs.Q1[Position](world)
	fmt.Println("Positions:")
	for q.Next() {
		_, pos := q.Item()
		fmt.Println(pos)
	}
	ecs.AddStartupSystem(world, func(c *ecs.Commands) {
		ecs.Spawn2(c, Position{X: 500, Y: 1000}, Velocity{X: 100, Y: -100})
		q := ecs.Q1[Position](c.World())
		for q.Next() {
			e, pos := q.Item()
			fmt.Println(pos)
			if pos.X == 7 {
				ecs.RemoveComponent[Velocity](c, e)
			}
		}
	})
	world.Step()
	q.Reset()
	fmt.Println("Positions (second pass):")
	for q.Next() {
		_, pos := q.Item()
		fmt.Println(pos)
	}
	world.Step()
	q.Reset()
	fmt.Println("Positions (third pass):")
	for q.Next() {
		_, pos := q.Item()
		fmt.Println(pos)
	}
}

type Position struct {
	X, Y float64
}

// implement ecs.Component interface
func (Position) ComponentUUID() string {
	return "558ae276-f21d-4251-94d5-0f0b3941f420"
}

type Velocity struct {
	X, Y float64
}

// implement ecs.Component interface
func (Velocity) ComponentUUID() string {
	return "14805b66-ed17-49ad-9a1f-75589e8465a2"
}

type GlobalProps struct {
	ecs.Resource
	TotalPosVelChanged int
}

// implement ecs.Resource interface
func (GlobalProps) ResourceUUID() string {
	return "9e1ba9f4-6211-4b79-b02a-41aaec79a10c"
}

func moveAll(commands *ecs.Commands) {
	// props.TotalPosVelChanged = 0
	// for posVelQuery.Next() {
	// 	_, pos, vel := posVelQuery.Item()
	// 	pos.X += vel.X
	// 	pos.Y += vel.Y
	// 	props.TotalPosVelChanged++
	// }
	q := ecs.Q2[Position, Velocity](commands.World())
	for q.Next() {
		_, pos, vel := q.Item()
		pos.X += vel.X
		pos.Y += vel.Y
	}
}

func setup(commands *ecs.Commands) {
	ecs.Spawn2(commands, Position{X: 5, Y: 10}, Velocity{X: 1, Y: 2})
}
