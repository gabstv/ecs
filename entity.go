package ecs

// Entity is a unique identifier. It is used to group components in a World.
type Entity uint64

// EntityEvent is used for event handlers of entities, views and systems
type EntityEvent func(e Entity, w *World)
