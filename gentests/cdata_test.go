package gentests

import (
	"testing"
	"unsafe"

	"github.com/gabstv/ecs/v3"
)

func TestResize(t *testing.T) {
	w := ecs.NewWorld()
	ecs.RegisterWorldDefaults(w)

	e1 := w.NewEntity()
	SetPositionComponentData(w, e1, Position{
		X: 10,
		Y: 10,
	})
	SetRotationComponentData(w, e1, Rotation{
		Angle: 45,
	})
	wr1 := WatchRotationComponentData(w, e1)
	p00 := uintptr(unsafe.Pointer(wr1.Data()))
	//
	e2 := w.NewEntity()
	SetPositionComponentData(w, e2, Position{
		X: 20,
		Y: 40,
	})
	SetRotationComponentData(w, e2, Rotation{
		Angle: 90,
	})
	//
	wr1.Data().Angle++
	//
	e3 := w.NewEntity()
	SetPositionComponentData(w, e3, Position{
		X: 100,
		Y: -100,
	})
	SetRotationComponentData(w, e3, Rotation{
		Angle: 180,
	})
	//
	p11 := uintptr(unsafe.Pointer(wr1.Data()))
	//
	if p00 == p11 {
		t.Fatalf("p00 %v == p11 %v", p00, p11)
	}
	if wr1.Data().Angle != 46 {
		t.Fatalf("angle %v != 46", wr1.Data().Angle)
	}
}

func TestEvents(t *testing.T) {
	w := ecs.NewWorld()
	ecs.RegisterWorldDefaults(w)

	x := 0
	y := 0

	exid := w.Listen(ecs.EvtComponentAdded|ecs.EvtComponentRemoved, func(e ecs.Event) {
		x++
	})
	eyid := w.Listen(ecs.EvtComponentRemoved, func(e ecs.Event) {
		y++
	})

	e1 := w.NewEntity()
	SetPositionComponentData(w, e1, Position{
		X: 10,
		Y: 10,
	})
	GetPositionComponent(w).Remove(e1)

	if x != 2 {
		t.FailNow()
	}
	if y != 1 {
		t.FailNow()
	}

	w.RemoveListener(eyid)

	SetPositionComponentData(w, e1, Position{
		X: 12,
		Y: -16,
	})

	w.RemoveListener(exid)

	if x != 3 {
		t.FailNow()
	}
}
