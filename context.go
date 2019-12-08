package ecs

import (
	"context"
	"time"
)

type Context interface {
	context.Context
	DT() float64
	System() *System
}

type ctxt struct {
	c      context.Context
	dt     float64
	system *System
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

func (c ctxt) WithSystem(s *System) Context {
	clone := ctxt{
		c:      c.c,
		dt:     c.dt,
		system: s,
	}
	return clone
}
