package ecs

import (
	"fmt"
	"math"
	"testing"
)

type testPosition struct {
	X float64
	Y float64
}

// benchmarks

func BenchmarkAddEntity(b *testing.B) {
	w := NewWorld()
	for i := 0; i < b.N; i++ {
		w.NewEntity()
	}
}

func BenchmarkAdd10Entities(b *testing.B) {
	w := NewWorld()
	for i := 0; i < b.N; i += 10 {
		w.NewEntities(10)
	}
}

func BenchmarkAdd100Entities(b *testing.B) {
	w := NewWorld()
	for i := 0; i < b.N; i += 100 {
		w.NewEntities(100)
	}
}

func BenchmarkAdd1000Entities(b *testing.B) {
	w := NewWorld()
	for i := 0; i < b.N; i += 1000 {
		w.NewEntities(1000)
	}
}

func BenchmarkPointerComponent(b *testing.B) {
	b.StopTimer()
	w := NewWorld()
	comp, _ := w.NewComponent(NewComponentInput{
		Name: "position",
	})
	view := w.NewView(comp)
	for i := 0; i < 1000; i++ {
		entity := w.NewEntity()
		w.AddComponentToEntity(entity, comp, &testPosition{})
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		matches := view.Matches()
		for _, mm := range matches {
			ev := mm.Components[comp].(*testPosition)
			ev.X += 0.01
			ev.Y += 0.02
		}
	}
}

func BenchmarkRun10000PositionX2(b *testing.B) {
	b.StopTimer()
	w := NewWorld()
	comp, _ := w.NewComponent(NewComponentInput{
		Name: "position",
	})
	w.NewSystem(0, func(dt float64, view *View) {
		ls := view.Matches()
		for _, v := range ls {
			j := v.Components[comp].(*testPosition)
			j.X += 0.05
		}
	}, comp)
	w.NewSystem(-1, func(dt float64, view *View) {
		ls := view.Matches()
		for _, v := range ls {
			j := v.Components[comp].(*testPosition)
			j.Y += 0.03
		}
	}, comp)
	for i := 0; i < 10000; i++ {
		entity := w.NewEntity()
		w.AddComponentToEntity(entity, comp, &testPosition{})
	}
	b.StartTimer()
	// run!
	for i := 0; i < b.N; i++ {
		w.Run(1.0)
	}
}

func BenchmarkRun10000PositionX2Tagged(b *testing.B) {
	b.StopTimer()
	w := NewWorld()
	comp, _ := w.NewComponent(NewComponentInput{
		Name: "position",
	})
	upd := w.NewSystem(0, func(dt float64, view *View) {
		ls := view.Matches()
		for _, v := range ls {
			j := v.Components[comp].(*testPosition)
			j.X += 0.05
		}
	}, comp)
	upd.AddTag("update")
	drw := w.NewSystem(-1, func(dt float64, view *View) {
		ls := view.Matches()
		for _, v := range ls {
			j := v.Components[comp].(*testPosition)
			j.Y += 0.03
		}
	}, comp)
	drw.AddTag("draw")
	for i := 0; i < 10000; i++ {
		entity := w.NewEntity()
		w.AddComponentToEntity(entity, comp, &testPosition{})
	}
	b.StartTimer()
	// run!
	for i := 0; i < b.N; i++ {
		w.RunWithTag("update", 1.0)
		w.RunWithTag("draw", 1.0)
	}
}

func BenchmarkRun10000PositionX3(b *testing.B) {
	b.StopTimer()
	w := NewWorld()
	comp, _ := w.NewComponent(NewComponentInput{
		Name: "position",
	})
	w.NewSystem(0, func(dt float64, view *View) {
		ls := view.Matches()
		for _, v := range ls {
			j := v.Components[comp].(*testPosition)
			j.X += 0.05
		}
	}, comp)
	w.NewSystem(-1, func(dt float64, view *View) {
		ls := view.Matches()
		for _, v := range ls {
			j := v.Components[comp].(*testPosition)
			j.Y += 0.03
		}
	}, comp)
	w.NewSystem(-2, func(dt float64, view *View) {
		ls := view.Matches()
		for _, v := range ls {
			j := v.Components[comp].(*testPosition)
			j.Y += math.Sqrt(j.X + 1)
		}
	}, comp)
	for i := 0; i < 10000; i++ {
		entity := w.NewEntity()
		w.AddComponentToEntity(entity, comp, &testPosition{})
	}
	b.StartTimer()
	// run!
	for i := 0; i < b.N; i++ {
		w.Run(1.0)
	}
}

// tests

func TestComponent(t *testing.T) {
	w := NewWorld()
	entity := w.NewEntity()
	comp, err := w.NewComponent(NewComponentInput{
		Name: "position",
	})
	if err != nil {
		t.Fatal(err)
	}
	if comp == nil {
		t.Fatal("comp is null")
	}
	err = w.AddComponentToEntity(entity, comp, &testPosition{2, 3})
	if err != nil {
		t.Fatal(err)
	}
	res := w.Query(comp)
	if len(res) != 1 {
		t.Fatal("query should return a result")
	}
	dd := res[0].Components[comp].(*testPosition)
	dd.X += 0.1
	dd.Y += 0.1
	res = w.Query(comp)
	if len(res) != 1 {
		t.Fatal("query should return a result")
	}
	dd = res[0].Components[comp].(*testPosition)
	if dd.X != 2.1 {
		t.Fatal("dd.X != 2.1")
	}
	if dd.Y != 3.1 {
		t.Fatal("dd.Y != 3.1")
	}
	err = w.RemoveComponentFromEntity(entity, comp)
	if err != nil {
		t.Fatal(err)
	}
	res = w.Query(comp)
	if len(res) != 0 {
		t.Fatal("query should return no results")
	}
}

func TestSystemSort(t *testing.T) {
	w := NewWorld()
	comp, _ := w.NewComponent(NewComponentInput{
		Name: "position",
	})
	w.NewSystem(1, func(dt float64, view *View) {
		fmt.Println("sys1", view.Matches())
	}, comp)
	w.NewSystem(100, func(dt float64, view *View) {
		fmt.Println("sys100", view.Matches())
	}, comp)
	w.NewSystem(-1000, func(dt float64, view *View) {
		fmt.Println("syslast", view.Matches())
	}, comp)
	if w.systems[0].priority != 100 {
		t.Fail()
	}
	if w.systems[2].priority != -1000 {
		t.Fail()
	}
}

func TestRun(t *testing.T) {
	w := NewWorld()
	comp, _ := w.NewComponent(NewComponentInput{
		Name: "position",
	})
	var lastEntityPos *testPosition
	w.NewSystem(-1000, func(dt float64, view *View) {
		ls := view.Matches()
		lastEntityPos = ls[9999].Components[comp].(*testPosition)
	}, comp)
	w.NewSystem(0, func(dt float64, view *View) {
		ls := view.Matches()
		for _, v := range ls {
			j := v.Components[comp].(*testPosition)
			j.X++
			j.Y += 2
		}
	}, comp)
	w.NewSystem(-1, func(dt float64, view *View) {
		ls := view.Matches()
		for _, v := range ls {
			j := v.Components[comp].(*testPosition)
			j.X *= 4
			j.Y *= 4
		}
	}, comp)
	for i := 0; i < 10000; i++ {
		entity := w.NewEntity()
		w.AddComponentToEntity(entity, comp, &testPosition{})
	}
	// run!
	w.Run(1.0)
	if lastEntityPos == nil {
		t.Fail()
	}
	if lastEntityPos.X != 4 {
		t.Fail()
	}
	if lastEntityPos.Y != 8 {
		t.Fail()
	}
}
