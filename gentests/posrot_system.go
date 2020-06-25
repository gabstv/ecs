// Code generated by ecs https://github.com/gabstv/ecs; DO NOT EDIT.

package gentests

import (
    
    "sort"

    "github.com/gabstv/ecs/v2"
    
)









const uuidPosRotSystem = "58FFC3BE-7BC8-4381-A93B-74945405F171"

type viewPosRotSystem struct {
    entities []VIPosRotSystem
    world ecs.BaseWorld
    
}

type VIPosRotSystem struct {
    Entity ecs.Entity
    
    Position *Position 
    
    Rotation *Rotation 
    
}

type sortedVIPosRotSystems []VIPosRotSystem
func (a sortedVIPosRotSystems) Len() int           { return len(a) }
func (a sortedVIPosRotSystems) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortedVIPosRotSystems) Less(i, j int) bool { return a[i].Entity < a[j].Entity }

func newviewPosRotSystem(w ecs.BaseWorld) *viewPosRotSystem {
    return &viewPosRotSystem{
        entities: make([]VIPosRotSystem, 0),
        world: w,
    }
}

func (v *viewPosRotSystem) Matches() []VIPosRotSystem {
    
    return v.entities
    
}

func (v *viewPosRotSystem) indexof(e ecs.Entity) int {
    i := sort.Search(len(v.entities), func(i int) bool { return v.entities[i].Entity >= e })
    if i < len(v.entities) && v.entities[i].Entity == e {
        return i
    }
    return -1
}

// Fetch a specific entity
func (v *viewPosRotSystem) Fetch(e ecs.Entity) (data VIPosRotSystem, ok bool) {
    
    i := v.indexof(e)
    if i == -1 {
        return VIPosRotSystem{}, false
    }
    return v.entities[i], true
}

func (v *viewPosRotSystem) Add(e ecs.Entity) bool {
    
    
    // MUST NOT add an Entity twice:
    if i := v.indexof(e); i > -1 {
        return false
    }
    v.entities = append(v.entities, VIPosRotSystem{
        Entity: e,
        Position: GetPositionComponent(v.world).Data(e),
Rotation: GetRotationComponent(v.world).Data(e),

    })
    if len(v.entities) > 1 {
        if v.entities[len(v.entities)-1].Entity < v.entities[len(v.entities)-2].Entity {
            sort.Sort(sortedVIPosRotSystems(v.entities))
        }
    }
    return true
}

func (v *viewPosRotSystem) Remove(e ecs.Entity) bool {
    
    
    if i := v.indexof(e); i != -1 {

        v.entities = append(v.entities[:i], v.entities[i+1:]...)
        return true
    }
    return false
}

func (v *viewPosRotSystem) clearpointers() {
    
    
    for i := range v.entities {
        e := v.entities[i].Entity
        
        v.entities[i].Position = nil
        
        v.entities[i].Rotation = nil
        
        _ = e
    }
}

func (v *viewPosRotSystem) rescan() {
    
    
    for i := range v.entities {
        e := v.entities[i].Entity
        
        v.entities[i].Position = GetPositionComponent(v.world).Data(e)
        
        v.entities[i].Rotation = GetRotationComponent(v.world).Data(e)
        
        _ = e
        
    }
}

// PosRotSystem implements ecs.BaseSystem
type PosRotSystem struct {
    initialized bool
    world       ecs.BaseWorld
    view        *viewPosRotSystem
    enabled     bool
    
}

// GetPosRotSystem returns the instance of the system in a World
func GetPosRotSystem(w ecs.BaseWorld) *PosRotSystem {
    return w.S(uuidPosRotSystem).(*PosRotSystem)
}

// Enable system
func (s *PosRotSystem) Enable() {
    s.enabled = true
}

// Disable system
func (s *PosRotSystem) Disable() {
    s.enabled = false
}

// Enabled checks if enabled
func (s *PosRotSystem) Enabled() bool {
    return s.enabled
}

// UUID implements ecs.BaseSystem
func (PosRotSystem) UUID() string {
    return "58FFC3BE-7BC8-4381-A93B-74945405F171"
}

func (PosRotSystem) Name() string {
    return "PosRotSystem"
}

// ensure matchfn
var _ ecs.MatchFn = matchPosRotSystem

// ensure resizematchfn
var _ ecs.MatchFn = resizematchPosRotSystem

func (s *PosRotSystem) match(eflag ecs.Flag) bool {
    return matchPosRotSystem(eflag, s.world)
}

func (s *PosRotSystem) resizematch(eflag ecs.Flag) bool {
    return resizematchPosRotSystem(eflag, s.world)
}

func (s *PosRotSystem) ComponentAdded(e ecs.Entity, eflag ecs.Flag) {
    if s.match(eflag) {
        if s.view.Add(e) {
            // TODO: dispatch event that this entity was added to this system
            
        }
    } else {
        if s.view.Remove(e) {
            // TODO: dispatch event that this entity was removed from this system
            
        }
    }
}

func (s *PosRotSystem) ComponentRemoved(e ecs.Entity, eflag ecs.Flag) {
    if s.match(eflag) {
        if s.view.Add(e) {
            // TODO: dispatch event that this entity was added to this system
            
        }
    } else {
        if s.view.Remove(e) {
            // TODO: dispatch event that this entity was removed from this system
            
        }
    }
}

func (s *PosRotSystem) ComponentResized(cflag ecs.Flag) {
    if s.resizematch(cflag) {
        s.view.rescan()
        
    }
}

func (s *PosRotSystem) ComponentWillResize(cflag ecs.Flag) {
    if s.resizematch(cflag) {
        
        s.view.clearpointers()
    }
}

func (s *PosRotSystem) V() *viewPosRotSystem {
    return s.view
}

func (*PosRotSystem) Priority() int64 {
    return 0
}

func (s *PosRotSystem) Setup(w ecs.BaseWorld) {
    if s.initialized {
        panic("PosRotSystem called Setup() more than once")
    }
    s.view = newviewPosRotSystem(w)
    s.world = w
    s.enabled = true
    s.initialized = true
    
}


func init() {
    ecs.RegisterSystem(func() ecs.BaseSystem {
        return &PosRotSystem{}
    })
}
