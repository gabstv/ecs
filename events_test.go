package ecs

import (
	"testing"
)

func TestComponentAddedEvent(t *testing.T) {

	type testEvcomp struct {
		Name string
	}

	w := NewWorld()

	nadded := 0
	addedNames := make([]string, 0)

	AddSystem(w, func(ctx *Context) {
		cpadded := ComponentsAdded[testEvcomp](ctx)
		for {
			cpair, ok := cpadded()
			if !ok {
				break
			}
			nadded++
			addedNames = append(addedNames, cpair.ComponentCopy.Name)
		}
	})
	AddStartupSystem(w, func(ctx *Context) {
		Spawn(ctx, testEvcomp{
			Name: "bacon",
		})
		Spawn(ctx, testEvcomp{
			Name: "pizza",
		})
		Spawn(ctx, testEvcomp{
			Name: "burger",
		})
	})
	w.Step()
	if nadded != 3 {
		t.Fatal("nadded != 3")
	}
	if len(addedNames) != 3 {
		t.Fatal("len(addedNames) != 3")
	}
	if addedNames[0] != "bacon" {
		t.Fatal("addedNames[0] != bacon")
	}
}
