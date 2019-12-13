package ecs

import (
	"context"
	"testing"
)

func TestSystemMiddleware(t *testing.T) {
	x := 100
	ffn := SysWrapFn(func(ctx Context) {
		x += ctx.System().Get("x").(int)
	}, func(next SystemExec) SystemExec {
		return func(ctx Context) {
			ctx.System().Set("x", 1000)
			next(ctx)
		}
	})
	ffn(ctxt{
		c:  context.Background(),
		dt: 0,
		system: &System{
			dict: newdict(),
		},
	})
	if x != 1100 {
		t.Fatal(x)
	}
}
