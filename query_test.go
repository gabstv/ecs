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

func TestSanity(t *testing.T) {
	type a struct {
		name int
	}
	type b struct {
		name int
	}
	ta := reflect.TypeOf(a{1})
	tb := reflect.TypeOf(b{1})
	if ta == tb {
		t.Error("types are equal")
	}
	za := zeroValue(ta)
	zb := zeroValue(tb)
	if za == zb {
		t.Error("zero values are equal")
	}
}
