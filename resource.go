package ecs

import "reflect"

// InitResource creates a new resource instance inside the world.
// The resource must be a struct. It panics if the resource already exists, or
// if you pass a pointer.
// If the struct implements DefaultResource, then InitResource will call Init.
//
// Example:
//
//	type MyResource struct {
//		// ...
//	}
//
//	func (r *MyResource) Init(w World) {
//		// ...
//	}
func InitResource[T any](w World) {
	var zt T
	t := reflect.TypeOf(zt)
	if t.Kind() == reflect.Ptr {
		panic("InitResource: resource must be a struct, not a pointer")
	}
	w.setResource(typeMapKeyOf(t), &zt)
}

// Resource retrieves a previously registered resource.
func Resource[T any](ctx *Context) *T {
	var zt T
	x := ctx.world.getResource(typeMapKeyOf(reflect.TypeOf(zt)))
	if x == nil {
		return nil
	}
	return x.(*T)
}

type WorldIniter interface {
	Init(w World)
}
