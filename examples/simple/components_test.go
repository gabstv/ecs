package simple

import (
	"math"
	"testing"
)

func TestECS(t *testing.T) {
	w := NewWorld()
	e := w.NewEntity()
	SetPositionComponentData(w, e, Position{
		X: 10,
		Y: 15,
	})
	SetRotationComponentData(w, e, Rotation{
		Radians: math.Pi,
	})
	SetVelocityComponentData(w, e, Velocity{
		X:       1,
		Y:       2,
		Radians: 0,
	})
	w.Update(1. / 60)
	if GetPositionComponentData(w, e).X != 11 {
		t.Fatal()
	}
	if GetPositionComponentData(w, e).Y != 17 {
		t.Fatal()
	}
	// remove Velocity{} so the Movement component wont have this entity anymore
	GetVelocityComponent(w).Remove(e)
	w.Update(1. / 60)
	if GetPositionComponentData(w, e).X != 11 {
		t.Fatal()
	}
	if GetPositionComponentData(w, e).Y != 17 {
		t.Fatal()
	}
}
