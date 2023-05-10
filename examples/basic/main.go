package main

import (
	"fmt"

	"github.com/gabstv/ecs/v4"
)

func main() {
	world := ecs.NewWorld()
	// the system below will run on every world.Step() call
	ecs.AddSystem(world, moveAll)
	// startup systems are executed only once, and before any other system
	ecs.AddStartupSystem(world, setup)
	world.Step()
	world.Step()
	q := ecs.Q1[Position](world)
	fmt.Println("Positions:")
	for q.Next() {
		_, pos := q.Item()
		fmt.Println(pos)
	}
	ecs.AddStartupSystem(world, func(c *ecs.Context) {
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
	props := ecs.GetResource[GlobalProps](world)
	fmt.Println("TotalPosVelChanged:", props.TotalPosVelChanged)
}

type Position struct {
	X, Y float64
}

type Velocity struct {
	X, Y float64
}

type GlobalProps struct {
	TotalPosVelChanged int
}

func moveAll(ctx *ecs.Context) {
	props := ecs.GetResource[GlobalProps](ctx.World())
	props.TotalPosVelChanged = 0
	q := ecs.Q2[Position, Velocity](ctx.World())
	for q.Next() {
		_, pos, vel := q.Item()
		pos.X += vel.X
		pos.Y += vel.Y
		props.TotalPosVelChanged++
	}
}

func setup(ctx *ecs.Context) {
	ecs.InitResource[GlobalProps](ctx.World())
	ecs.Spawn2(ctx, Position{X: 5, Y: 10}, Velocity{X: 1, Y: 2})
}
