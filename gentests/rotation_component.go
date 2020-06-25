// Code generated by ecs https://github.com/gabstv/ecs; DO NOT EDIT.

package gentests

import (
    "sort"
    

    "github.com/gabstv/ecs/v2"
)








const uuidRotationComponent = "56890133-3769-477A-B163-412C5ECC6B07"
const capRotationComponent = 2

type drawerRotationComponent struct {
    Entity ecs.Entity
    Data   Rotation
}

// WatchRotation is a helper struct to access a valid pointer of Rotation
type WatchRotation interface {
    Entity() ecs.Entity
    Data() *Rotation
}

type slcdrawerRotationComponent []drawerRotationComponent
func (a slcdrawerRotationComponent) Len() int           { return len(a) }
func (a slcdrawerRotationComponent) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a slcdrawerRotationComponent) Less(i, j int) bool { return a[i].Entity < a[j].Entity }


type mWatchRotation struct {
    c *RotationComponent
    entity ecs.Entity
}

func (w *mWatchRotation) Entity() ecs.Entity {
    return w.entity
}

func (w *mWatchRotation) Data() *Rotation {
    
    
    id := w.c.indexof(w.entity)
    if id == -1 {
        return nil
    }
    return &w.c.data[id].Data
}

// RotationComponent implements ecs.BaseComponent
type RotationComponent struct {
    initialized bool
    flag        ecs.Flag
    world       ecs.BaseWorld
    wkey        [4]byte
    data        []drawerRotationComponent
    
}

// GetRotationComponent returns the instance of the component in a World
func GetRotationComponent(w ecs.BaseWorld) *RotationComponent {
    return w.C(uuidRotationComponent).(*RotationComponent)
}

// SetRotationComponentData updates/adds a Rotation to Entity e
func SetRotationComponentData(w ecs.BaseWorld, e ecs.Entity, data Rotation) {
    GetRotationComponent(w).Upsert(e, data)
}

// GetRotationComponentData gets the *Rotation of Entity e
func GetRotationComponentData(w ecs.BaseWorld, e ecs.Entity) *Rotation {
    return GetRotationComponent(w).Data(e)
}

// WatchRotationComponentData gets a pointer getter of an entity's Rotation.
//
// The pointer must not be stored because it may become invalid overtime.
func WatchRotationComponentData(w ecs.BaseWorld, e ecs.Entity) WatchRotation {
    return &mWatchRotation{
        c: GetRotationComponent(w),
        entity: e,
    }
}

// UUID implements ecs.BaseComponent
func (RotationComponent) UUID() string {
    return "56890133-3769-477A-B163-412C5ECC6B07"
}

// Name implements ecs.BaseComponent
func (RotationComponent) Name() string {
    return "RotationComponent"
}

func (c *RotationComponent) indexof(e ecs.Entity) int {
    i := sort.Search(len(c.data), func(i int) bool { return c.data[i].Entity >= e })
    if i < len(c.data) && c.data[i].Entity == e {
        return i
    }
    return -1
}

// Upsert creates or updates a component data of an entity.
// Not recommended to be used directly. Use SetRotationComponentData to change component
// data outside of a system loop.
func (c *RotationComponent) Upsert(e ecs.Entity, data interface{}) {
    v, ok := data.(Rotation)
    if !ok {
        panic("data must be Rotation")
    }
    
    id := c.indexof(e)
    
    if id > -1 {
        
        dwr := &c.data[id]
        dwr.Data = v
        
        return
    }
    
    rsz := false
    if cap(c.data) == len(c.data) {
        rsz = true
        c.world.CWillResize(c, c.wkey)
        
    }
    newindex := len(c.data)
    c.data = append(c.data, drawerRotationComponent{
        Entity: e,
        Data:   v,
    })
    if len(c.data) > 1 {
        if c.data[newindex].Entity < c.data[newindex-1].Entity {
            c.world.CWillResize(c, c.wkey)
            
            sort.Sort(slcdrawerRotationComponent(c.data))
            rsz = true
        }
    }
    
    if rsz {
        
        c.world.CResized(c, c.wkey)
        c.world.Dispatch(ecs.Event{
            Type: ecs.EvtComponentsResized,
            ComponentName: "RotationComponent",
            ComponentID: "56890133-3769-477A-B163-412C5ECC6B07",
        })
    }
    
    c.world.CAdded(e, c, c.wkey)
    c.world.Dispatch(ecs.Event{
        Type: ecs.EvtComponentAdded,
        ComponentName: "RotationComponent",
        ComponentID: "56890133-3769-477A-B163-412C5ECC6B07",
        Entity: e,
    })
}

// Remove a Rotation data from entity e
//
// Warning: DO NOT call remove inside the system entities loop
func (c *RotationComponent) Remove(e ecs.Entity) {
    
    
    i := c.indexof(e)
    if i == -1 {
        return
    }
    
    //c.data = append(c.data[:i], c.data[i+1:]...)
    c.data = c.data[:i+copy(c.data[i:], c.data[i+1:])]
    c.world.CRemoved(e, c, c.wkey)
    
    c.world.Dispatch(ecs.Event{
        Type: ecs.EvtComponentRemoved,
        ComponentName: "RotationComponent",
        ComponentID: "56890133-3769-477A-B163-412C5ECC6B07",
        Entity: e,
    })
}

func (c *RotationComponent) Data(e ecs.Entity) *Rotation {
    
    
    index := c.indexof(e)
    if index > -1 {
        return &c.data[index].Data
    }
    return nil
}

// Flag returns the 
func (c *RotationComponent) Flag() ecs.Flag {
    return c.flag
}

// Setup is called by ecs.BaseWorld
//
// Do not call this directly
func (c *RotationComponent) Setup(w ecs.BaseWorld, f ecs.Flag, key [4]byte) {
    if c.initialized {
        panic("RotationComponent called Setup() more than once")
    }
    c.flag = f
    c.world = w
    c.wkey = key
    c.data = make([]drawerRotationComponent, 0, 2)
    c.initialized = true
}


func init() {
    ecs.RegisterComponent(func() ecs.BaseComponent {
        return &RotationComponent{}
    })
}
