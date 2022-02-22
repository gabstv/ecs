package ecs

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"

	"github.com/BurntSushi/toml"
	"github.com/gabstv/container"
)

// ComponentType is a data type that has a Pkg() function.
// It is used to identify a component type.
// It must be unique per type in a world.
type ComponentType interface {
	Pkg() string
}

// ComponentData holds the data T of an Entity.
type ComponentData[T ComponentType] struct {
	Entity Entity
	Data   T
}

// IComponentStore is an interface for component stores.
type IComponentStore interface {
	Contains(e Entity) bool
	Remove(e Entity) bool
	MergeJSONData(e Entity, jd []byte) error

	dataExtract(fn func(e Entity, d interface{}))
	dataImport(e Entity, d toml.Primitive, md toml.MetaData) error
	dataOf(e Entity) interface{}
	typeMatch(d interface{}) bool
}

// ComponentStore[T ComponentType] is a component data storage. The component data
// is stored in a slice ordered by the Entity (ID; ascending).
type ComponentStore[T ComponentType] struct {
	// data is an ordered entity slice of all entities that have this component.
	data  []ComponentData[T]
	world *World
	zerov T
	isptr *bool

	watchers container.Set[*ComponentWatcher[T]]
}

// Apply passes a pointer of the component data to the function fn.
// This is used to read or update data in the component.
func (c *ComponentStore[T]) Apply(e Entity, fn func(*T)) bool {
	index, exists := c.getIndex(e)
	if !exists {
		return false
	}
	x := &c.data[index]
	fn(&x.Data)
	return true
}

// Contains returns true if the entity has data of this component store.
func (c *ComponentStore[T]) Contains(e Entity) bool {
	_, exists := c.getIndex(e)
	return exists
}

// MergeJSONData unmarshals the data into the component type of this component store.
func (c *ComponentStore[T]) MergeJSONData(e Entity, jd []byte) error {
	zv, _ := c.getCopy(e)
	if err := json.Unmarshal(jd, &zv); err != nil {
		return err
	}
	c.Replace(e, zv)
	return nil
}

// Remove removes the component data from this component store. It returns true
// if the component data was found (and then removed).
func (c *ComponentStore[T]) Remove(e Entity) bool {
	index, exists := c.getIndex(e)
	if !exists {
		return false
	}
	c.data = append(c.data[:index], c.data[index+1:]...)
	c.watchers.Each(func(w *ComponentWatcher[T]) {
		w.ComponentRemoved(e)
	})
	return true
}

// Replace adds or replaces the component data for the given entity.
func (c *ComponentStore[T]) Replace(e Entity, data T) {
	// add to c.data
	index, exists := c.getIndex(e)
	if exists {
		c.setDataAt(index, data)
		return
	}
	// insert data at index
	c.data = Insert(c.data, index, ComponentData[T]{e, data})
	c.watchers.Each(func(w *ComponentWatcher[T]) {
		w.ComponentAdded(e)
	})
}

// ComponentStore[T] privates

func (c *ComponentStore[T]) all() []ComponentData[T] {
	return c.data
}

func (c *ComponentStore[T]) dataExtract(fn func(e Entity, d interface{})) {
	for _, v := range c.all() {
		fn(v.Entity, v.Data)
	}
}

func (c *ComponentStore[T]) dataImport(e Entity, d toml.Primitive, md toml.MetaData) error {
	x := c.newType()
	if c.isPointerType() {
		if err := md.PrimitiveDecode(d, x); err != nil {
			return fmt.Errorf("failed to decode component %T: %v", x, err)
		}
	} else {
		if err := md.PrimitiveDecode(d, &x); err != nil {
			return fmt.Errorf("failed to decode component %T: %v", x, err)
		}
	}
	c.Replace(e, x)
	return nil
}

func (c *ComponentStore[T]) dataOf(e Entity) interface{} {
	i, exists := c.getIndex(e)
	if !exists {
		return nil
	}
	return c.data[i].Data
}

func (c *ComponentStore[T]) getCopy(e Entity) (T, bool) {
	index, exists := c.getIndex(e)
	if !exists {
		return c.zerov, false
	}
	return c.data[index].Data, true
}

// getIndex does a binary search on c.data for the index of the entity
func (c *ComponentStore[T]) getIndex(e Entity) (int, bool) {
	return getIndex(c.data, e)
}

func (c *ComponentStore[T]) isPointerType() bool {
	if c.isptr == nil {
		//TODO: consider dropping support of pointer types to avoid using
		//      reflect entirely.
		if reflect.ValueOf(c.zerov).Kind() == reflect.Ptr {
			isptr := true
			c.isptr = &isptr
		} else {
			isptr := false
			c.isptr = &isptr
		}
	}
	return *c.isptr
}

func (c *ComponentStore[T]) newType() T {
	if c.isPointerType() {
		//TODO: consider dropping support of pointer types to avoid using
		//      reflect entirely.
		return reflect.New(reflect.TypeOf(c.zerov)).Interface().(T)
	}
	var y T
	return y
}

func (c *ComponentStore[T]) setDataAt(index int, data T) {
	c.data[index].Data = data
}

func (c *ComponentStore[T]) typeMatch(d interface{}) bool {
	if d == nil {
		return false
	}
	_, ok := d.(T)
	return ok
}

type ComponentIndexEntry struct {
	Name  string
	Index int
}

type ComponentIndex []ComponentIndexEntry

func (e ComponentIndex) ToMap() map[string]int {
	m := make(map[string]int)
	for _, v := range e {
		m[v.Name] = v.Index
	}
	return m
}

// static fns

// Apply updates the component data for the given entity.
func Apply[T ComponentType](w *World, e Entity, fn func(*T)) bool {
	c := GetComponentStore[T](w)
	return c.Apply(e, fn)
}

// Contains returns true if the given entity has the given component.
func Contains[T ComponentType](w *World, e Entity) bool {
	c := GetComponentStore[T](w)
	return c.Contains(e)
}

// GetComponentStore returns the component store for the given component type and
// world instance.
func GetComponentStore[T ComponentType](w *World) *ComponentStore[T] {
	if w.components == nil {
		w.components = make(map[string]IComponentStore)
	}
	var zv T
	if c, ok := w.components[zv.Pkg()]; ok {
		return c.(*ComponentStore[T])
	}
	c := &ComponentStore[T]{
		data:  make([]ComponentData[T], 0),
		world: w,
	}
	w.components[zv.Pkg()] = c
	return c
}

// RemoveComponent removes the component data for the given entity.
// It returns false if the component was not found.
func RemoveComponent[T ComponentType](w *World, e Entity) bool {
	c := GetComponentStore[T](w)
	return c.Remove(e)
}

// Set replaces or inserts the component data for the given entity.
// If you junst need to update a value, use Apply() instead.
func Set[T ComponentType](w *World, e Entity, data T) {
	c := GetComponentStore[T](w)
	c.Replace(e, data)
}

func componentIndexFromMap(m map[string]int) ComponentIndex {
	x := make([]ComponentIndexEntry, 0, len(m))
	for k, v := range m {
		x = append(x, ComponentIndexEntry{
			Name:  k,
			Index: v,
		})
	}
	sort.SliceStable(x, func(i, j int) bool {
		return x[i].Index < x[j].Index
	})
	return ComponentIndex(x)
}

// getEntityIndex does a binary search on an entity slice
func getEntityIndex(slc []Entity, e Entity) (int, bool) {
	x := sort.Search(len(slc), func(i int) bool {
		return slc[i] >= e
	})
	if x < len(slc) && slc[x] == e {
		// x is present at data[i]
		return x, true
	}
	// x is not present in data,
	// but i is the index where it would be inserted.
	return x, false
}

// getIndex does a binary search on c.data for the index of the entity
func getIndex[T ComponentType](slc []ComponentData[T], e Entity) (int, bool) {
	x := sort.Search(len(slc), func(i int) bool {
		return slc[i].Entity >= e
	})
	if x < len(slc) && slc[x].Entity == e {
		// x is present at data[i]
		return x, true
	}
	// x is not present in data,
	// but i is the index where it would be inserted.
	return x, false
}
