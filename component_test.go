package ecs

import "testing"

type TestMarco struct {
	X float64
	Y float64
	Z float64
}

type TestPolo struct {
	Scale float64
}

func BenchmarkViews100000(b *testing.B) {
	b.StopTimer()
	w := NewWorld()
	c1, _ := w.NewComponent(NewComponentInput{
		Name: "cmarco",
	})
	c2, _ := w.NewComponent(NewComponentInput{
		Name: "cpolo",
	})
	w.NewSystem("smarco", 100, func(ctx Context) {
		matches := ctx.System().View().Matches()
		for _, v := range matches {
			a := v.Components[c1].(*TestMarco)
			a.X += 0.1
			a.Y += 0.2
			a.Z += 0.3
		}
	}, c1)
	fmax := func(a, b, c float64) float64 {
		if a > b && a > c {
			return a
		}
		if b > c {
			return b
		}
		return c
	}
	w.NewSystem("smarcopolo", 99, func(ctx Context) {
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
