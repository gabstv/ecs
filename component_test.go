package ecs

import "testing"

type Position struct {
	X, Y int
}

func (Position) Pkg() string {
	return "test.Position"
}

type Rotation struct {
	Value int
}

func (Rotation) Pkg() string {
	return "test.Rotation"
}

func TestComponents(t *testing.T) {
	w := NewWorld()
	e0 := w.NewEntity()
	e1 := w.NewEntity()
	Set(w, e0, Position{
		X: 1,
		Y: 2,
	})
	Set(w, e1, Position{
		X: 3,
		Y: 4,
	})
	Set(w, e0, Rotation{
		Value: 1,
	})
	var p0 Position
	var p1 Position

	ok0 := Apply(w, e0, func(v *Position) {
		p0 = *v
	})
	ok1 := Apply(w, e1, func(v *Position) {
		p1 = *v
	})
	if !ok0 {
		t.Errorf("Get[Position] of e0 failed")
	}
	if !ok1 {
		t.Errorf("Get[Position] of e1 failed")
	}
	if p0.X != 1 {
		t.Errorf("Get[Position] of e0.X failed")
	}
	if p0.Y != 2 {
		t.Errorf("Get[Position] of e0.Y failed")
	}
	if p1.X != 3 {
		t.Errorf("Get[Position] of e1.X failed")
	}
	if p1.Y != 4 {
		t.Errorf("Get[Position] of e1.Y failed")
	}
	var r0 Rotation
	ok0 = Apply(w, e0, func(v *Rotation) {
		r0 = *v
	})
	var r1 Rotation
	ok1 = Apply(w, e1, func(v *Rotation) {
		r1 = *v
	})
	if !ok0 {
		t.Errorf("Get[Rotation] of e0 failed")
	}
	if ok1 {
		t.Errorf("Get[Rotation] of e1 should fail")
	}
	if r0.Value != 1 {
		t.Errorf("Get[Rotation] of e0.Value failed")
	}
	if r1.Value != 0 {
		t.Errorf("Get[Rotation] of e1.Value should fail")
	}
	if !RemoveComponent[Rotation](w, e0) {
		t.Errorf("Remove[Rotation] of e0 failed")
	}
	ok0 = Contains[Rotation](w, e0)
	if ok0 {
		t.Errorf("Get[Rotation] of e0 should fail")
	}
}
