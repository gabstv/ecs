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
	w.NewSystem("", 0, func(ctx Context) {
		ls := ctx.System().View().Matches()
		for _, v := range ls {
			j := v.Components[comp].(*testPosition)
			j.X += 0.05
		}
	}, comp)
	w.NewSystem("", -1, func(ctx Context) {
		ls := ctx.System().View().Matches()
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
	upd := w.NewSystem("", 0, func(ctx Context) {
		ls := ctx.System().View().Matches()
		for _, v := range ls {
			j := v.Components[comp].(*testPosition)
			j.X += 0.05
		}
	}, comp)
	upd.AddTag("update")
	drw := w.NewSystem("", -1, func(ctx Context) {
		ls := ctx.System().View().Matches()
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
	w.NewSystem("", 0, func(ctx Context) {
		ls := ctx.System().View().Matches()
		for _, v := range ls {
			j := v.Components[comp].(*testPosition)
			j.X += 0.05
		}
	}, comp)
	w.NewSystem("", -1, func(ctx Context) {
		ls := ctx.System().View().Matches()
		for _, v := range ls {
			j := v.Components[comp].(*testPosition)
			j.Y += 0.03
		}
	}, comp)
	w.NewSystem("", -2, func(ctx Context) {
		ls := ctx.System().View().Matches()
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

func TestQueryMask(t *testing.T) {
	w := NewWorld()
	xentity := w.NewEntity()
	xcomp, err := w.NewComponent(NewComponentInput{
		Name: "position",
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	if err := w.AddComponentToEntity(xentity, xcomp, &testPosition{X: 10, Y: 20}); err != nil {
		t.Fatal(err.Error())
	}
	yentity := w.NewEntity()
	ycomp, err := w.NewComponent(NewComponentInput{
		Name: "examplemask",
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	if err := w.AddComponentToEntity(yentity, xcomp, &testPosition{X: 11, Y: 21}); err != nil {
		t.Fatal(err.Error())
	}
	if err := w.AddComponentToEntity(yentity, ycomp, &testPosition{X: -1, Y: -1}); err != nil {
		t.Fatal(err.Error())
	}
	res := w.QueryMask(Cs(ycomp), Cs(xcomp))
	if len(res) != 1 {
		t.Fatal("invalid length")
	}
	if res[0].Components[xcomp] == nil {
		t.Fatal("invalid result")
	}
	if res[0].Components[xcomp].(*testPosition).X != 10 {
		t.Fatal("invalid result (X)")
	}
	if res[0].Components[ycomp] != nil {
		t.Fatal("invalid result (should be nil)")
	}
}

func TestMaskView(t *testing.T) {
	w := NewWorld()
	xentity := w.NewEntity()
	xcomp, err := w.NewComponent(NewComponentInput{
		Name: "position",
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	if err := w.AddComponentToEntity(xentity, xcomp, &testPosition{X: 10, Y: 20}); err != nil {
		t.Fatal(err.Error())
	}
	yentity := w.NewEntity()
	ycomp, err := w.NewComponent(NewComponentInput{
		Name: "examplemask",
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	zcomp, err := w.NewComponent(NewComponentInput{
		Name: "examplemask2",
	})
	if err != nil {
		t.Fatal(err.Error())
	}
	if err := w.AddComponentToEntity(yentity, xcomp, &testPosition{X: 11, Y: 21}); err != nil {
		t.Fatal(err.Error())
	}
	if err := w.AddComponentToEntity(yentity, ycomp, &testPosition{X: -1, Y: -1}); err != nil {
		t.Fatal(err.Error())
	}
	if err := w.AddComponentToEntity(yentity, zcomp, &testPosition{X: -2, Y: -2}); err != nil {
		t.Fatal(err.Error())
	}
	view := w.NewMaskView(Cs(ycomp, zcomp), Cs(xcomp))
	view.SetOnEntityAdded(func(e Entity, w *World) {
		ydat := xcomp.Data(e)
		if ydat == nil {
			t.Fatal("ydat is nil (view event)")
		}
		ydat2, ok := ydat.(*testPosition)
		if !ok {
			t.Fatal("ydat is not *testPosition (view event)")
		}
		if ydat2.X != 11 {
			t.Fatal("ydat is not correct (view event)")
		}
	})
	{
		m := view.Matches()
		if len(m) != 1 {
			t.Fatal("invalid length")
		}
		if m[0].Components[xcomp] == nil {
			t.Fatal("invalid result")
		}
		if m[0].Components[xcomp].(*testPosition).X != 10 {
			t.Fatal("invalid result (X)")
		}
		if m[0].Components[ycomp] != nil {
			t.Fatal("invalid result (should be nil)")
		}
	}
	if err := w.RemoveComponentFromEntity(yentity, ycomp); err != nil {
		t.Fatal(err.Error())
	}
	{
		m := view.Matches()
		if len(m) != 1 {
			t.Fatal("invalid length")
		}
	}
	if err := w.RemoveComponentFromEntity(yentity, zcomp); err != nil {
		t.Fatal(err.Error())
	}
	{
		m := view.Matches()
		if len(m) != 2 {
			t.Fatal("invalid length")
		}
	}
}

func TestSystemSort(t *testing.T) {
	w := NewWorld()
	comp, _ := w.NewComponent(NewComponentInput{
		Name: "position",
	})
	w.NewSystem("", 1, func(ctx Context) {
		fmt.Println("sys1", ctx.System().View().Matches())
	}, comp)
	w.NewSystem("", 100, func(ctx Context) {
		fmt.Println("sys100", ctx.System().View().Matches())
	}, comp)
	w.NewSystem("", -1000, func(ctx Context) {
		fmt.Println("syslast", ctx.System().View().Matches())
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
	w.NewSystem("", -1000, func(ctx Context) {
		ls := ctx.System().View().Matches()
		lastEntityPos = ls[9999].Components[comp].(*testPosition)
	}, comp)
	w.NewSystem("", 0, func(ctx Context) {
		ls := ctx.System().View().Matches()
		for _, v := range ls {
			j := v.Components[comp].(*testPosition)
			j.X++
			j.Y += 2
		}
	}, comp)
	w.NewSystem("", -1, func(ctx Context) {
		ls := ctx.System().View().Matches()
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

func TestDict(t *testing.T) {
	w := NewWorld()
	a := w.Get("a")
	if a != nil {
		t.Fail()
	}
	w.Set("a", int64(19))
	a64 := w.Get("a").(int64)
	if a64 != 19 {
		t.Fail()
	}
}

func TestSystemDict(t *testing.T) {
	w := NewWorld()
	comp, _ := w.NewComponent(NewComponentInput{
		Name: "comp",
	})
	sys := w.NewSystem("", 0, func(ctx Context) {
		ctx.System().Set("marco", "polo")
	}, comp)
	entity0 := w.NewEntity()
	w.AddComponentToEntity(entity0, comp, &testPosition{})
	w.Run(1)
	poloiface := sys.Get("marco")
	if poloiface == nil {
		t.FailNow()
	}
	if str, ok := poloiface.(string); ok {
		if str != "polo" {
			t.Fail()
		}
	} else {
		t.Fail()
	}
}
