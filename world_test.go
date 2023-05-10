package ecs

import (
	"math/rand"
	"reflect"
	"testing"
)

func BenchmarkReflectionTechniques(b *testing.B) {
	type basics struct {
		foo int
	}
	for i := 0; i < b.N; i++ {
		cmd4 := reflect.TypeOf(&basics{
			foo: i * 2,
		})
		_ = cmd4
	}
}

func Benchmark10kEntitiesSimple(b *testing.B) {
	b.StopTimer()
	type Position struct {
		X, Y, Z float64
	}
	type Velocity struct {
		X, Y, Z float64
	}
	w := NewWorld()
	AddSystem(w, func(c *Context) {
		pvquery := Q2[Position, Velocity](c.World())
		for pvquery.Next() {
			_, pos, vel := pvquery.Item()
			pos.X += vel.X
			pos.Y += vel.Y
			pos.Z += vel.Z
		}
	})
	AddStartupSystem(w, func(c *Context) {
		for i := 0; i < 10000; i++ {
			Spawn2(c, Position{
				X: rand.Float64(),
				Y: rand.Float64(),
				Z: rand.Float64(),
			},
				Velocity{
					X: rand.Float64(),
					Y: rand.Float64(),
					Z: rand.Float64(),
				})
		}
	})
	w.Step()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		w.Step()
	}
}

func Benchmark10kEntitiesSimple5050(b *testing.B) {
	b.StopTimer()
	type Position struct {
		X, Y, Z float64
	}
	type Velocity struct {
		X, Y, Z float64
	}
	w := NewWorld()
	AddSystem(w, func(c *Context) {
		pvquery := Q2[Position, Velocity](c.World())
		for pvquery.Next() {
			_, pos, vel := pvquery.Item()
			pos.X += vel.X
			pos.Y += vel.Y
			pos.Z += vel.Z
		}
	})
	AddStartupSystem(w, func(c *Context) {
		for i := 0; i < 10000; i++ {
			if i%3 == 0 {
				Spawn2(c, Position{
					X: rand.Float64(),
					Y: rand.Float64(),
					Z: rand.Float64(),
				},
					Velocity{
						X: rand.Float64(),
						Y: rand.Float64(),
						Z: rand.Float64(),
					})
			} else if i%2 == 0 {
				Spawn(c, Velocity{
					X: rand.Float64(),
					Y: rand.Float64(),
					Z: rand.Float64(),
				})
			} else {
				Spawn(c, Position{
					X: rand.Float64(),
					Y: rand.Float64(),
					Z: rand.Float64(),
				})
			}
		}
	})
	w.Step()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		w.Step()
	}
}

func Benchmark100kEntitiesSimple(b *testing.B) {
	b.StopTimer()
	type Position struct {
		X, Y, Z float64
	}
	type Velocity struct {
		X, Y, Z float64
	}
	w := NewWorld()
	AddSystem(w, func(c *Context) {
		pvquery := Q2[Position, Velocity](c.World())
		for pvquery.Next() {
			_, pos, vel := pvquery.Item()
			pos.X += vel.X
			pos.Y += vel.Y
			pos.Z += vel.Z
		}
	})
	AddStartupSystem(w, func(c *Context) {
		for i := 0; i < 100000; i++ {
			Spawn2(c, Position{
				X: rand.Float64(),
				Y: rand.Float64(),
				Z: rand.Float64(),
			},
				Velocity{
					X: rand.Float64(),
					Y: rand.Float64(),
					Z: rand.Float64(),
				})
		}
	})
	w.Step()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		w.Step()
	}
}

func Benchmark1MEntitiesSimple(b *testing.B) {
	b.StopTimer()
	type Position struct {
		X, Y, Z float64
	}
	type Velocity struct {
		X, Y, Z float64
	}
	w := NewWorld()
	AddSystem(w, func(c *Context) {
		pvquery := Q2[Position, Velocity](c.World())
		for pvquery.Next() {
			_, pos, vel := pvquery.Item()
			pos.X += vel.X
			pos.Y += vel.Y
			pos.Z += vel.Z
		}
	})
	AddStartupSystem(w, func(c *Context) {
		for i := 0; i < 1000000; i++ {
			Spawn2(c, Position{
				X: rand.Float64(),
				Y: rand.Float64(),
				Z: rand.Float64(),
			},
				Velocity{
					X: rand.Float64(),
					Y: rand.Float64(),
					Z: rand.Float64(),
				})
		}
	})
	w.Step()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		w.Step()
	}
}

func TestRemoveEntity(t *testing.T) {
	type Person struct {
		Name string
	}
	type Education struct {
		Degree float64
	}
	w := NewWorld()
	AddStartupSystem(w, func(c *Context) {
		Spawn2(c, Person{Name: "Bob"}, Education{Degree: 3.0})
		Spawn2(c, Person{Name: "Alice"}, Education{Degree: 4.0})
	})
	AddSystem(w, func(c *Context) {
		query := Q2[Person, Education](c.World())
		for query.Next() {
			entt, _, e := query.Item()
			if e.Degree < 3.1 {
				RemoveEntity(c, entt)
			}
		}
	})
	w.Step()
	query := Q2[Person, Education](w)
	for query.Next() {
		_, p, _ := query.Item()
		if p.Name == "Bob" {
			t.Error("Bob should be removed")
		}
	}
	// add again
	AddStartupSystem(w, func(c *Context) {
		Spawn2(c, Person{Name: "Roomba"}, Education{Degree: 5.0})
	})
	w.Step()
	roombaEnt := Entity(0)
	// remove the person component
	AddStartupSystem(w, func(c *Context) {
		query := Q1[Person](c.World())
		for query.Next() {
			entt, p := query.Item()
			if p.Name == "Roomba" {
				roombaEnt = entt
				RemoveComponent[Person](c, entt)
			}
		}
	})
	w.Step()
	query = Q2[Person, Education](w)
	for query.Next() {
		_, p, _ := query.Item()
		if p.Name == "Roomba" {
			t.Error("Roomba should be removed")
		}
	}
	// add roomba again
	AddStartupSystem(w, func(c *Context) {
		qr := Q1[Education](c.World())
		for qr.Next() {
			entt, e := qr.Item()
			if e.Degree == 5.0 {
				AddComponent(c, entt, Person{Name: "Roomba 2"})
			}
		}
	})
	w.Step()
	query = Q2[Person, Education](w)
	foundRoomba2 := false
	for query.Next() {
		entt, p, _ := query.Item()
		if p.Name == "Roomba 2" {
			foundRoomba2 = true
			if entt != roombaEnt {
				t.Error("Roomba 2 should be the same entity")
			}
		}
	}
	if !foundRoomba2 {
		t.Error("Roomba 2 not found")
	}
}
