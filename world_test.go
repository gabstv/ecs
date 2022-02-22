package ecs

import (
	"bytes"
	_ "embed"
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

type BenchPos3 struct {
	X, Y, Z float64
}

func (BenchPos3) Pkg() string {
	return "test.BenchPos3"
}

type BenchSpeed3 struct {
	Xs, Ys, Zs float64
}

func (BenchSpeed3) Pkg() string {
	return "test.BenchSpeed3"
}

type BenchDeltaSpeed struct {
	Delta float64
}

func (BenchDeltaSpeed) Pkg() string {
	return "test.BenchDeltaSpeed"
}

type BenchAccel struct {
	Xa, Ya, Za float64
}

func (BenchAccel) Pkg() string {
	return "test.BenchAccel"
}

func BenchmarkWorldRuntimeWith4Systems(b *testing.B) {
	b.StopTimer()
	w := NewWorld()
	for i := 0; i < 2000; i++ {
		e := w.NewEntity()
		Set(w, e, BenchPos3{
			X: rand.Float64() * 10,
			Y: rand.Float64() * 10,
			Z: rand.Float64() * 10,
		})
		if i%2 == 0 {
			Set(w, e, BenchSpeed3{
				Xs: rand.Float64(),
				Ys: rand.Float64(),
				Zs: rand.Float64(),
			})
		}
		if i%3 == 0 {
			Set(w, e, BenchDeltaSpeed{
				Delta: 0,
			})
		}
	}
	updatePosition := NewSystem2[BenchPos3, BenchSpeed3](-100, w)
	updatePosition.Run = func(view *View2[BenchPos3, BenchSpeed3]) {
		view.Each(func(e Entity, d1 *BenchPos3, d2 *BenchSpeed3) {
			d1.X += d2.Xs * 1.0 / 60.0
			d1.Y += d2.Ys * 1.0 / 60.0
			d1.Z += d2.Zs * 1.0 / 60.0
		})
	}
	updateSpeed := NewSystem2[BenchSpeed3, BenchAccel](-101, w)
	updateSpeed.Run = func(view *View2[BenchSpeed3, BenchAccel]) {
		view.Each(func(e Entity, d1 *BenchSpeed3, d2 *BenchAccel) {
			d1.Xs += d2.Xa * 1.0 / 60.0
			d1.Ys += d2.Ya * 1.0 / 60.0
			d1.Zs += d2.Za * 1.0 / 60.0
		})
	}
	updAccelDamping := NewSystem[BenchAccel](-99, w)
	updAccelDamping.Run = func(view *View[BenchAccel]) {
		view.Each(func(e Entity, d *BenchAccel) {
			d.Xa *= 0.99
			d.Ya *= 0.99
			d.Za *= 0.99
		})
	}
	updateDeltaSpeed := NewSystem2[BenchSpeed3, BenchDeltaSpeed](-98, w)
	updateDeltaSpeed.Run = func(view *View2[BenchSpeed3, BenchDeltaSpeed]) {
		view.Each(func(_ Entity, d1 *BenchSpeed3, d2 *BenchDeltaSpeed) {
			d2.Delta = (d1.Xs + d1.Ys + d1.Zs) / 3
		})
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		w.Step()
	}
	// display last  for debug purposes
	slc := GetComponentStore[BenchPos3](w).data
	b.Log(slc[len(slc)-1])
}

//go:embed world_test_m.txt
var expectedText string

func TestMarshaler(t *testing.T) {

	w := NewWorld()
	e := w.NewEntity()
	idh := w.EntityUUID(e).String()
	Set(w, e, BenchPos3{
		X: 1,
		Y: 2,
		Z: 3,
	})
	Set(w, e, BenchSpeed3{
		Xs: 4,
		Ys: 5,
		Zs: 6,
	})
	Set(w, e, BenchDeltaSpeed{
		Delta: 0,
	})
	Set(w, e, BenchAccel{
		Xa: 0.1,
		Ya: 0.2,
		Za: 0.3,
	})
	buf := new(bytes.Buffer)
	assert.NoError(t, w.MarshalTo(buf))
	actualv := buf.String()
	assert.Equal(t, fmt.Sprintf(expectedText, idh), actualv)
}

func TestUnmarshaler(t *testing.T) {
	w := NewWorld()
	// ensure that components are registered
	_ = GetComponentStore[BenchPos3](w)
	_ = GetComponentStore[BenchSpeed3](w)
	_ = GetComponentStore[BenchDeltaSpeed](w)
	_ = GetComponentStore[BenchAccel](w)
	buf := bytes.NewBufferString(fmt.Sprintf(expectedText, "1B20DDAE-FE41-41D6-BC7F-EE46C175ED32"))
	assert.NoError(t, w.UnmarshalFrom(buf))
	assert.Equal(t, 1, len(w.entities))
	var pos BenchPos3
	ok := GetComponentStore[BenchPos3](w).Apply(w.entities[0], func(p *BenchPos3) {
		pos = *p
	})
	assert.True(t, ok)
	assert.Equal(t, 2.0, pos.Y)
}

func TestWorldRemoveEntity(t *testing.T) {
	w := NewWorld()
	e := w.NewEntity()
	Set(w, e, BenchPos3{
		X: 10.0,
		Y: 10.1,
		Z: 10.3,
	})
	Set(w, e, BenchSpeed3{
		Xs: .1,
		Ys: .2,
		Zs: .3,
	})
	sys := NewSystem2[BenchPos3, BenchSpeed3](0, w)
	sys.Run = func(view *View2[BenchPos3, BenchSpeed3]) {
		view.Each(func(e Entity, p *BenchPos3, s *BenchSpeed3) {
			p.X += s.Xs
			p.Y += s.Ys
			p.Z += s.Zs
		})
	}
	w.Step()
	Apply(w, e, func(p *BenchPos3) {
		assert.Less(t, 10.09, p.X)
		assert.Less(t, 10.29, p.Y)
		assert.Less(t, 10.59, p.Z)
	})
	assert.True(t, Remove(w, e))
	assert.False(t, Remove(w, e))
	assert.False(t, Apply(w, e, func(p *BenchPos3) {}))
}
