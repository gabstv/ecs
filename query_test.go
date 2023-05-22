package ecs

import (
	"reflect"
	"testing"
)

func isPointerType[T any]() bool {
	var zv T
	t := reflect.TypeOf(zv)
	return t.Kind() == reflect.Ptr
}

func TestQueryBasics(t *testing.T) {
	type Bacon struct {
	}
	if isPointerType[Bacon]() {
		t.Error("Bacon is not a pointer type")
	}
	if !isPointerType[*Bacon]() {
		t.Error("Bacon is a pointer type")
	}
}
