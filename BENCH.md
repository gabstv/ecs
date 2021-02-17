# V1 -> V2 Comparison

This is an archive of the difference between the code generated, pure ECS of
the new version vs the interface based V1 of ECS (1.x)

```
BenchmarkV1Views100000-4   	     175	   6904582 ns/op	      96 B/op	       2 allocs/op
BenchmarkV2Views100000-4   	    1780	    682237 ns/op	      32 B/op	       1 allocs/op
```

With the same components and sytems, v2.x is 10x faster.

bench_test.go
```go
package bench

import (
	"testing"

	"github.com/gabstv/ecs/v1/ecs"
)

func BenchmarkV1Views100000(b *testing.B) {
	b.StopTimer()
	w := ecs.NewWorld()
	c1, _ := w.NewComponent(ecs.NewComponentInput{
		Name: "cmarco",
	})
	c2, _ := w.NewComponent(ecs.NewComponentInput{
		Name: "cpolo",
	})
	w.NewSystem("smarco", 100, func(ctx ecs.Context) {
		matches := ctx.System().View().Matches()
		for _, v := range matches {
			a := v.Components[c1].(*TestMarco)
			a.X += 0.1
			a.Y += 0.2
			a.Z += 0.3
		}
	}, c1)
	w.NewSystem("smarcopolo", 99, func(ctx ecs.Context) {
		matches := ctx.System().View().Matches()
		for _, v := range matches {
			a := v.Components[c1].(*TestMarco)
			b := v.Components[c2].(*TestPolo)
			b.Scale = fmax(a.X, a.Y, a.Z)
			v.Components[c2] = b
		}
	}, c1, c2)
	entities := w.NewEntities(100000)
	for i, e := range entities {
		w.AddComponentToEntity(e, c1, &TestMarco{})
		if i%7 == 0 {
			w.AddComponentToEntity(e, c2, &TestPolo{})
		}
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		w.Run(1 / 60)
	}
}

func BenchmarkV2Views100000(b *testing.B) {
	b.StopTimer()
	w := NewWorld()
	for i := 0; i < 100000; i++ {
		e := w.NewEntity()
		SetTestMarcoComponentData(w, e, TestMarco{})
		if i%7 == 0 {
			SetTestPoloComponentData(w, e, TestPolo{})
		}
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		w.Update(1 / 60)
	}
}
```

bench.go
```go
package bench

import "github.com/gabstv/ecs"

//go:generate go run ../../../cmd/ecsgen/main.go -n TestMarco -p bench -o testmarco_comp.go -t ../../../templates/component.tmpl --vars "UUID=131E1425-DF74-4F50-9521-54845314E1EA"

type TestMarco struct {
	X float64 // 8 bytes
	Y float64 // 8 bytes
	Z float64 // 8 bytes
}

//go:generate go run ../../../cmd/ecsgen/main.go -n TestPolo -p bench -o testpolo_comp.go -t ../../../templates/component.tmpl --vars "UUID=C0EA39E3-9EB1-4415-92A5-3DB2AF0D6750"

type TestPolo struct {
	Scale float64 // 8 bytes
}

type System interface {
	ecs.System
	Update(dt float64)
}

type World struct {
	*ecs.World
}

func (w *World) Update(dt float64) {
	w.EachSystem(func(s ecs.System) bool {
		s.(System).Update(dt)
		return true
	})
}

func NewWorld() *World {
	w := &World{
		World: &ecs.World{},
	}
	w.World.Init()
	ecs.RegisterWorldDefaults(w)
	return w
}

//go:generate go run ../../../cmd/ecsgen/main.go -n Marco -p bench -o marco_system.go -t ../../../templates/system.tmpl --vars "UUID=A8C56294-1FE1-4BE8-BABF-29697B822D1B" --vars "Priority=100" --components "TestMarco"
//go:generate go run ../../../cmd/ecsgen/main.go -n MarcoPolo -p bench -o marcopolo_system.go -t ../../../templates/system.tmpl --vars "UUID=7EA57F45-51A4-4367-B147-340A279C091A" --vars "Priority=99" --components "TestMarco" --components "TestPolo"

var matchMarcoSystem = func(f ecs.Flag, w ecs.World) bool {
	return f.Contains(GetTestMarcoComponent(w).Flag())
}

var matchMarcoPoloSystem = func(f ecs.Flag, w ecs.World) bool {
	return f.Contains(GetTestMarcoComponent(w).Flag().Or(GetTestPoloComponent(w).Flag()))
}

func (s *MarcoSystem) Update(dt float64) {
	for _, x := range s.V().Matches() {
		x.TestMarco.X += .1
		x.TestMarco.Y += .2
		x.TestMarco.Z += .3
	}
}

func fmax(a, b, c float64) float64 {
	if a > b && a > c {
		return a
	}
	if b > c {
		return b
	}
	return c
}

func (s *MarcoPoloSystem) Update(dt float64) {
	for _, e := range s.V().Matches() {
		m := e.TestMarco
		p := e.TestPolo
		p.Scale = fmax(m.X, m.Y, m.Z)
	}
}

```