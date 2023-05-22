package ecs

import "reflect"

type ResourceUUID string

type ResourceT interface {
	ResourceUUID() ResourceUUID
}

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
func InitResource[T ResourceT](w World) {
	var zt T
	t := reflect.TypeOf(zt)
	if t.Kind() == reflect.Ptr {
		panic("InitResource: resource must be a struct, not a pointer")
	}
	w.setResource(zt.ResourceUUID(), &zt)
}

// Resource retrieves a previously registered resource.
func Resource[T ResourceT](ctx *Context) *T {
	var zt T
	x := ctx.world.getResource(zt.ResourceUUID())
	if x == nil {
		return nil
	}
	return x.(*T)
}

type WorldIniter interface {
	Init(w World)
}
