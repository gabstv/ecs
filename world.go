package ecs

import (
	"fmt"
	"io"
	"sort"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/google/uuid"
)

type Printer interface {
	Printf(fmt string, args ...interface{})
}

var (
	SerializerLogger Printer
)

type World struct {
	lastEntity  Entity
	entityMutex sync.Mutex
	entities    []Entity
	entityIDs   map[Entity]uuid.UUID // this is used when serializing/deserializing data
	entityUUIDs map[uuid.UUID]Entity
	components  map[string]IComponentStore
	systems     []ISystem
	sysMap      map[int]ISystem
	sysid       int
	isloading   bool
	enabled     bool
}

func NewWorld() *World {
	return &World{
		lastEntity:  0,
		entities:    make([]Entity, 0, 1024),
		entityIDs:   make(map[Entity]uuid.UUID),
		entityUUIDs: make(map[uuid.UUID]Entity),
		components:  make(map[string]IComponentStore),
		systems:     make([]ISystem, 0, 32),
		sysMap:      make(map[int]ISystem),
		enabled:     true,
	}
}

func (w *World) IsLoading() bool {
	return w.isloading
}

// EntityUUID returns the UUID of the entity
// If the entity exists, but no UUID is set, a new UUID is generated and set
func (w *World) EntityUUID(e Entity) uuid.UUID {
	x := sort.Search(len(w.entities), func(i int) bool {
		return w.entities[i] >= e
	})
	if x < len(w.entities) && w.entities[x] == e {
		// x is present at data[i]
		if uuid, ok := w.entityIDs[e]; ok {
			return uuid
		}
		id, err := uuid.NewRandom()
		if err != nil {
			panic(err)
		}
		w.entityIDs[e] = id
		w.entityUUIDs[id] = e
		return id
	}
	// entity not found
	return uuid.UUID{}
}

// getEntityByUUID is used by the deserializer to get the entity with the given UUID
//
// If the entity does not exist, a new entity is created and is assigned with the UUID
func (w *World) getEntityByUUID(id uuid.UUID) Entity {
	if e, ok := w.entityUUIDs[id]; ok {
		return e
	}
	// create a new entity and set the uuid to it
	e := w.NewEntity()
	w.entityUUIDs[id] = e
	w.entityIDs[e] = id
	return e
}

// EntityByUUID returns the entity with the given UUID
// If the entity does not exist, an empty entity (0) is returned
func (w *World) EntityByUUID(id uuid.UUID) (Entity, bool) {
	if e, ok := w.entityUUIDs[id]; ok {
		return e, true
	}
	return 0, false
}

func (w *World) NewEntity() Entity {
	w.entityMutex.Lock()
	defer w.entityMutex.Unlock()
	w.lastEntity++
	w.entities = append(w.entities, w.lastEntity)
	return w.lastEntity
}

func (w *World) NewEntities(count int) []Entity {
	if count <= 0 {
		return nil
	}
	w.entityMutex.Lock()
	defer w.entityMutex.Unlock()
	entts := make([]Entity, count)
	for i := 0; i < count; i++ {
		entts[i] = w.lastEntity + Entity(i+1)
	}
	w.lastEntity += Entity(count)
	w.entities = append(w.entities, entts...)
	return entts
}

var DefaultWorld = NewWorld()

func (w *World) addSystem(sys ISystem) int {
	w.sysid++
	w.systems = append(w.systems, sys)

	sort.SliceStable(w.systems, func(i, j int) bool {
		return w.systems[i].Priority() < w.systems[j].Priority()
	})
	w.sysMap[w.sysid] = sys

	return w.sysid
}

func (w *World) RemoveSystem(id int) bool {
	if sys, ok := w.sysMap[id]; ok {
		delete(w.sysMap, id)
		di := -1
		for i, s := range w.systems {
			if s == sys {
				di = i
				break
			}
		}
		if di >= 0 {
			w.systems = append(w.systems[:di], w.systems[di+1:]...)
		}
		return true
	}
	return false
}

func (w *World) Step() {
	for _, sys := range w.systems {
		sys.Execute()
	}
}

func (w *World) StepF(flag int) {
	for _, sys := range w.systems {
		if sys.Flag()&flag != 0 {
			sys.Execute()
		}
	}
}

func (w *World) Enabled() bool {
	return w.enabled
}

func (w *World) SetEnabled(v bool) {
	w.enabled = v
}

// AllEntities returns all entities in the world
func (w *World) AllEntities() []Entity {
	w.entityMutex.Lock()
	defer w.entityMutex.Unlock()
	ecopy := make([]Entity, len(w.entities))
	copy(ecopy, w.entities)
	return ecopy
}

type SerializedWorld struct {
	Entities       []SerializedEntity `toml:"entities"`
	ComponentIndex ComponentIndex     `toml:"component_index"`
	Enabled        bool               `toml:"enabled"`
}

type DeserializedWorld struct {
	Entities       []DeserializedEntity `toml:"entities"`
	ComponentIndex ComponentIndex       `toml:"component_index"`
	Enabled        bool                 `toml:"enabled"`
}

type DeserializedEntity struct {
	UUID       uuid.UUID                   `toml:"uuid"`
	Components []DeserializedComponentData `toml:"components"`
}

type SerializedEntity struct {
	UUID       uuid.UUID     `toml:"uuid"`
	Components []interface{} `toml:"components"`
}

type DeserializedComponentData struct {
	CI   int            `toml:"ci"` // component index
	Data toml.Primitive `toml:"data"`
}

type SerializedComponentData struct {
	CI   int         `toml:"ci"` // component index
	Data interface{} `toml:"data"`
}

// MarshalTo marshals the world data to a writer
func (w *World) MarshalTo(dw io.Writer) error {
	return w.serializeData(toml.NewEncoder(dw))
}

func (w *World) serializeData(me Encoder) error {
	encoderMutex.Lock()
	defer encoderMutex.Unlock()
	setEncoderWorld(w)
	defer setEncoderWorld(nil)
	sw := SerializedWorld{
		Entities: make([]SerializedEntity, 0, len(w.entities)),
		Enabled:  w.enabled,
	}
	compIndex := make(map[string]int)
	entt := make(map[Entity]*SerializedEntity)

	type ctuple struct {
		Name  string
		Store IComponentStore
	}
	ctuples := make([]ctuple, 0, len(w.components))
	for name, c := range w.components {
		ctuples = append(ctuples, ctuple{Name: name, Store: c})
	}
	sort.SliceStable(ctuples, func(i, j int) bool {
		return ctuples[i].Name < ctuples[j].Name
	})

	for indexm, ctuplev := range ctuples {
		c := ctuplev.Store
		ci := indexm + 1
		compIndex[ctuplev.Name] = ci
		c.dataExtract(func(e Entity, d interface{}) {
			ent := entt[e]
			if ent == nil {
				ent = &SerializedEntity{
					UUID:       w.EntityUUID(e),
					Components: make([]interface{}, 0),
				}
			}
			ent.Components = append(ent.Components, SerializedComponentData{
				CI:   ci,
				Data: d,
			})
			entt[e] = ent
		})
	}
	stb := make([]Sortable[Entity, *SerializedEntity], 0, len(entt))
	for eid, ent := range entt {
		stb = append(stb, Sortable[Entity, *SerializedEntity]{
			Index: eid,
			Data:  ent,
		})
	}
	sort.Slice(stb, func(i, j int) bool {
		return stb[i].Index < stb[j].Index
	})
	for _, s := range stb {
		sw.Entities = append(sw.Entities, *s.Data)
	}
	sw.ComponentIndex = componentIndexFromMap(compIndex)
	return me.Encode(sw)
}

func (w *World) UnmarshalFrom(dr io.Reader) error {
	x := &DeserializedWorld{}
	md, err := toml.NewDecoder(dr).Decode(x)
	if err != nil {
		return fmt.Errorf("failed to decode toml world data: %w", err)
	}
	return w.deserializeData(md, x)
}

func (w *World) UnmarshalFromMeta(md toml.MetaData, prim toml.Primitive) error {
	x := &DeserializedWorld{}
	err := md.PrimitiveDecode(prim, x)
	if err != nil {
		return fmt.Errorf("failed to decode toml world data: %w", err)
	}
	return w.deserializeData(md, x)
}

func (w *World) GetGenericComponent(registryName string) IComponentStore {
	return w.components[registryName]
}

func (w *World) deserializeData(md toml.MetaData, dw *DeserializedWorld) error {
	w.isloading = true
	decoderMutex.Lock()
	defer decoderMutex.Unlock()
	setDecoderWorld(w)
	defer setDecoderWorld(nil)
	w.enabled = dw.Enabled
	compoSmap := dw.ComponentIndex.ToMap()
	compoImap := make(map[int]string)
	for k, v := range compoSmap {
		compoImap[v] = k
	}
	compos := make(map[int]IComponentStore)
	// the components need to be registered beforehand
	for i, v := range compoImap {
		compos[i] = w.GetGenericComponent(v)
	}
	for _, ent := range dw.Entities {
		e := w.getEntityByUUID(ent.UUID)
		for _, c := range ent.Components {
			if compos[c.CI] == nil {
				if SerializerLogger != nil {
					SerializerLogger.Printf("component [%d] %s not registered", c.CI, compoImap[c.CI])
				}
			} else {
				compos[c.CI].dataImport(e, c.Data, md)
			}
		}
	}
	w.isloading = false
	return nil
}

type Encoder interface {
	Encode(interface{}) error
}
