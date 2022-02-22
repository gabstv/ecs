package ecs

import (
	"fmt"
	"io"
	"sort"

	"github.com/BurntSushi/toml"
	"github.com/google/uuid"
)

type World struct {
	lastEntity   Entity
	entities     []Entity
	entityIDs    map[Entity]uuid.UUID // this is used when serializing/deserializing data
	entityUUIDs  map[uuid.UUID]Entity
	eventManager *eventManager
	components   map[string]IComponentStore
	systems      []ISystem
	sysMap       map[int]ISystem
	sysid        int
	isloading    bool
	enabled      bool
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
	w.lastEntity++
	w.entities = append(w.entities, w.lastEntity)
	return w.lastEntity
}

// Remove removes an Entity. It tries to delete the entity from all the
// component registries of this world.
func (w *World) Remove(e Entity) bool {
	x := sort.Search(len(w.entities), func(i int) bool {
		return w.entities[i] >= e
	})
	if x >= len(w.entities) || w.entities[x] != e {
		return false
	}
	// x is present at data[i]
	for _, c := range w.components {
		_ = c.Remove(e)
	}
	w.entities = append(w.entities[:x], w.entities[x+1:]...)
	return true
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
	ecopy := make([]Entity, len(w.entities))
	copy(ecopy, w.entities)
	return ecopy
}

// MarshalTo marshals the world data to a writer
func (w *World) MarshalTo(dw io.Writer) error {
	return w.serializeData(toml.NewEncoder(dw))
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

func (w *World) addSystem(sys ISystem) int {
	w.sysid++
	w.systems = append(w.systems, sys)

	sort.SliceStable(w.systems, func(i, j int) bool {
		return w.systems[i].Priority() < w.systems[j].Priority()
	})
	w.sysMap[w.sysid] = sys

	return w.sysid
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

// NewWorld creates a new world. A world is not thread safe .I t shouldn't be
// shared between threads. This function also adds the default systems to the
// world. To create a new world without any systems, use NewEmptyWorld()
func NewWorld() *World {
	w := newWorld()
	globalSystems.lock.Lock()
	defer globalSystems.lock.Unlock()
	for _, sf := range globalSystems.sysFactory {
		sf(w)
	}
	return w
}

// NewEmptyWorld creates a new world. A world is not thread safe .I t shouldn't
// be shared between threads. NewEmptyWorld creates a new world with no systems.
func NewEmptyWorld() *World {
	return newWorld()
}

// Remove removes an entity from the world. It will also remove all components
// attached to the entity. It returns false if the entity was not found.
func Remove(w *World, e Entity) bool {
	return w.Remove(e)
}

func newWorld() *World {
	return &World{
		lastEntity:   0,
		entities:     make([]Entity, 0, 1024),
		entityIDs:    make(map[Entity]uuid.UUID),
		entityUUIDs:  make(map[uuid.UUID]Entity),
		components:   make(map[string]IComponentStore),
		systems:      make([]ISystem, 0, 32),
		sysMap:       make(map[int]ISystem),
		enabled:      true,
		eventManager: newEventManager(),
	}
}
