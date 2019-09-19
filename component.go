package ecs

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// ValidateComponentData is a function that receives data and determines if
// the received data is the expected struct that the component should have.
type ValidateComponentData func(data interface{}) bool

// ComponentDestructor is an optional destructor function that a component might have.
type ComponentDestructor func(w *World, entity Entity, data interface{})

// Component is a container of raw data.
type Component struct {
	lock         sync.RWMutex
	data         map[Entity]interface{}
	destructor   ComponentDestructor
	name         string
	flag         flag
	validatedata ValidateComponentData
}

// String returns the string representation of this component
func (c *Component) String() string {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return fmt.Sprintf("[Component %q]", c.name)
}

// Validate if data belongs to the component.
func (c *Component) Validate(data interface{}) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if c.validatedata == nil {
		fmt.Printf("component %v called Validate but the validate func is nil\n", c.String())
		return false
	}
	return c.validatedata(data)
}

// Data returns the component data of an entity. Returns nil if the entity doesn't have this component.
func (c *Component) Data(e Entity) interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.data[e]
}

// NewComponentInput is the input args of the World.NewComponent function.
type NewComponentInput struct {
	// Name of the component. Used for debugging purposes.
	Name string
	// ValidateDataFn is the optional component data validator function.
	ValidateDataFn ValidateComponentData
	// DestructorFn is the optional destructor function of a component.
	DestructorFn ComponentDestructor
}

// NewComponent takes the input and returns a pointer of a component, used to
// add component data to an entity.
//
// This function is necessary because the data of a component doesn't have any
// metadata to identify which component it belongs to. That's because the
// data of a component is as simple as possible, without any embedded logic, just data.
//
// The max amount of component types, per system, is 256.
func (w *World) NewComponent(input NewComponentInput) (*Component, error) {
	nextid := atomic.AddUint64(&w.nextComponent, 1)
	if nextid > 256 {
		return nil, fmt.Errorf("component overflow (max = 256)")
	}
	comp := &Component{
		data:         make(map[Entity]interface{}),
		destructor:   input.DestructorFn,
		name:         input.Name,
		flag:         newflagbit(uint8(nextid - 1)), // index starts at 0
		validatedata: input.ValidateDataFn,
	}
	w.lock.Lock()
	w.components[comp.flag] = comp
	w.lock.Unlock()
	return comp, nil
}
