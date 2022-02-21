package ecs

import (
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/google/uuid"
)

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

type Encoder interface {
	Encode(interface{}) error
}

var encoderMutex sync.Mutex
var encoderWorld *World

var decoderMutex sync.Mutex
var decoderWorld *World

func setEncoderWorld(w *World) {
	encoderWorld = w
}

func setDecoderWorld(w *World) {
	decoderWorld = w
}
