// Code generated by ecs https://github.com/gabstv/ecs; DO NOT EDIT.

package gentests

import (
    "sort"
    

    "github.com/gabstv/ecs/v3"
)








const uuidPositionComponent = "3DF7F486-807D-4CE8-A187-37CED338137B"
const capPositionComponent = 256

type drawerPositionComponent struct {
    Entity ecs.Entity
    Data   Position
}

// WatchPosition is a helper struct to access a valid pointer of Position
type WatchPosition interface {
    Entity() ecs.Entity
    Data() *Position
}

type slcdrawerPositionComponent []drawerPositionComponent
func (a slcdrawerPositionComponent) Len() int           { return len(a) }
func (a slcdrawerPositionComponent) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a slcdrawerPositionComponent) Less(i, j int) bool { return a[i].Entity < a[j].Entity }


type mWatchPosition struct {
    c *PositionComponent
    entity ecs.Entity
}

func (w *mWatchPosition) Entity() ecs.Entity {
    return w.entity
}

func (w *mWatchPosition) Data() *Position {
    
    
    id := w.c.indexof(w.entity)
    if id == -1 {
        return nil
    }
    return &w.c.data[id].Data
}

// PositionComponent implements ecs.BaseComponent
type PositionComponent struct {
    initialized bool
    flag        ecs.Flag
    world       ecs.BaseWorld
    wkey        [4]byte
    data        []drawerPositionComponent
    
}

// GetPositionComponent returns the instance of the component in a World
func GetPositionComponent(w ecs.BaseWorld) *PositionComponent {
    return w.C(uuidPositionComponent).(*PositionComponent)
}

// SetPositionComponentData updates/adds a Position to Entity e
func SetPositionComponentData(w ecs.BaseWorld, e ecs.Entity, data Position) {
    GetPositionComponent(w).Upsert(e, data)
}

// GetPositionComponentData gets the *Position of Entity e
func GetPositionComponentData(w ecs.BaseWorld, e ecs.Entity) *Position {
    return GetPositionComponent(w).Data(e)
}

// WatchPositionComponentData gets a pointer getter of an entity's Position.
//
// The pointer must not be stored because it may become invalid overtime.
func WatchPositionComponentData(w ecs.BaseWorld, e ecs.Entity) WatchPosition {
    return &mWatchPosition{
        c: GetPositionComponent(w),
        entity: e,
    }
}

// UUID implements ecs.BaseComponent
func (PositionComponent) UUID() string {
    return "3DF7F486-807D-4CE8-A187-37CED338137B"
}

// Name implements ecs.BaseComponent
func (PositionComponent) Name() string {
    return "PositionComponent"
}

func (c *PositionComponent) indexof(e ecs.Entity) int {
    i := sort.Search(len(c.data), func(i int) bool { return c.data[i].Entity >= e })
    if i < len(c.data) && c.data[i].Entity == e {
        return i
    }
    return -1
}

// Upsert creates or updates a component data of an entity.
// Not recommended to be used directly. Use SetPositionComponentData to change component
// data outside of a system loop.
func (c *PositionComponent) Upsert(e ecs.Entity, data interface{}) {
    v, ok := data.(Position)
    if !ok {
        panic("data must be Position")
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
    c.data = append(c.data, drawerPositionComponent{
        Entity: e,
        Data:   v,
    })
    if len(c.data) > 1 {
        if c.data[newindex].Entity < c.data[newindex-1].Entity {
            c.world.CWillResize(c, c.wkey)
            
            sort.Sort(slcdrawerPositionComponent(c.data))
            rsz = true
        }
    }
    
    if rsz {
        
        c.world.CResized(c, c.wkey)
        c.world.Dispatch(ecs.Event{
            Type: ecs.EvtComponentsResized,
            ComponentName: "PositionComponent",
            ComponentID: "3DF7F486-807D-4CE8-A187-37CED338137B",
        })
    }
    
    c.world.CAdded(e, c, c.wkey)
    c.world.Dispatch(ecs.Event{
        Type: ecs.EvtComponentAdded,
        ComponentName: "PositionComponent",
        ComponentID: "3DF7F486-807D-4CE8-A187-37CED338137B",
        Entity: e,
    })
}

// Remove a Position data from entity e
//
// Warning: DO NOT call remove inside the system entities loop
func (c *PositionComponent) Remove(e ecs.Entity) {
    
    
    i := c.indexof(e)
    if i == -1 {
        return
    }
    
    //c.data = append(c.data[:i], c.data[i+1:]...)
    c.data = c.data[:i+copy(c.data[i:], c.data[i+1:])]
    c.world.CRemoved(e, c, c.wkey)
    
    c.world.Dispatch(ecs.Event{
        Type: ecs.EvtComponentRemoved,
        ComponentName: "PositionComponent",
        ComponentID: "3DF7F486-807D-4CE8-A187-37CED338137B",
        Entity: e,
    })
}

func (c *PositionComponent) Data(e ecs.Entity) *Position {
    
    
    index := c.indexof(e)
    if index > -1 {
        return &c.data[index].Data
    }
    return nil
}

// Flag returns the 
func (c *PositionComponent) Flag() ecs.Flag {
    return c.flag
}

// Setup is called by ecs.BaseWorld
//
// Do not call this directly
func (c *PositionComponent) Setup(w ecs.BaseWorld, f ecs.Flag, key [4]byte) {
    if c.initialized {
        panic("PositionComponent called Setup() more than once")
    }
    c.flag = f
    c.world = w
    c.wkey = key
    c.data = make([]drawerPositionComponent, 0, 256)
    c.initialized = true
}


func init() {
    ecs.RegisterComponent(func() ecs.BaseComponent {
        return &PositionComponent{}
    })
}
