package ecs

import (
	"context"
	"time"
)

type Context interface {
	context.Context
	DT() float64
	System() *System
	World() Worlder
}

type ContextBuilderFn func(c0 context.Context, dt float64, sys *System, w Worlder) Context

var DefaultContextBuilder = func(c0 context.Context, dt float64, sys *System, w Worlder) Context {
	return ctxt{
		c:      c0,
		dt:     dt,
		system: sys,
		world:  w,
	}
}

type ctxt struct {
	c      context.Context
	dt     float64
	system *System
	world  Worlder
}

func (c ctxt) Deadline() (deadline time.Time, ok bool) {
	return c.c.Deadline()
}

func (c ctxt) Done() <-chan struct{} {
	return c.c.Done()
}

func (c ctxt) Err() error {
	return c.c.Err()
}

func (c ctxt) Value(key interface{}) interface{} {
	return c.c.Value(key)
}

func (c ctxt) DT() float64 {
	return c.dt
}

func (c ctxt) System() *System {
	return c.system
}

func (c ctxt) World() Worlder {
	return c.world
}

func (c ctxt) WithSystem(s *System) Context {
	clone := ctxt{
		c:      c.c,
		dt:     c.dt,
		system: s,
	}
	return clone
}
