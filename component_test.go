package ecs

import "testing"

type position struct {
	x, y int
}

func (c position) ComponentUUID() ComponentUUID {
	return "test.position"
}

func TestWeakReference(t *testing.T) {

	type posDB struct {
		pcopy []WeakReference[position]
	}

	w := NewWorld()

	AddStartupSystem(w, func(ctx *Context) {
		Spawn(ctx, position{
			x: 1,
			y: 2,
		})
	})
	AddSystem(w, func(ctx *Context) {
		cpadded := ComponentsAdded[position](ctx)
		for {
			cpair, ok := cpadded()
			if !ok {
				break
			}
			wref := NewWeakReference[position](ctx, cpair.Entity)
			if wref == nil {
				t.Error("wref == nil")
			}
			db := LocalResource[posDB](ctx)
			db.pcopy = append(db.pcopy, wref)
		}
	})
	w.Step()
}
