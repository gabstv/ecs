package ecs

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/google/uuid"
)

type Entity uint64

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

func (e Entity) MarshalBinary() ([]byte, error) {
	id := encoderWorld.EntityUUID(e)
	if id.Version() < 1 {
		return nil, fmt.Errorf("entity %d has no UUID", e)
	}
	return id[:], nil
}

func (e *Entity) UnmarshalBinary(data []byte) error {
	if len(data) != 16 {
		return fmt.Errorf("invalid UUID length: %d", len(data))
	}
	var d [16]byte
	copy(d[:], data)
	id := uuid.UUID(d)

	if id.Version() < 1 {
		return fmt.Errorf("invalid UUID")
	}
	*e = decoderWorld.getEntityByUUID(id)
	return nil
}

func (e Entity) MarshalText() (text []byte, err error) {
	if e == 0 {
		return nil, nil
	}
	return []byte(encoderWorld.EntityUUID(e).String()), nil
}

func (e *Entity) UnmarshalText(text []byte) error {
	if string(text) == "" || string(text) == "00000000-0000-0000-0000-000000000000" {
		return nil
	}
	id, err := uuid.Parse(string(text))
	if err != nil {
		return err
	}
	*e = decoderWorld.getEntityByUUID(id)
	return nil
}

func (e *Entity) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if s == "" {
		return nil
	}
	if s == "00000000-0000-0000-0000-000000000000" {
		return nil
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return err
	}
	*e = decoderWorld.getEntityByUUID(id)
	return nil
}

func (e Entity) MarshalJSON() ([]byte, error) {
	s := encoderWorld.EntityUUID(e).String()
	return json.Marshal(s)
}

type Entities []Entity

func (e Entities) MarshalBinary() ([]byte, error) {
	eslice := make([]byte, len(e)*16+8)
	offset := binary.PutUvarint(eslice, uint64(len(e)))
	if offset != 8 {
		panic("unexpected offset")
	}
	for _, e := range e {
		id := encoderWorld.EntityUUID(e)
		if id.Version() < 1 {
			return nil, fmt.Errorf("entity %d has no UUID", e)
		}
		copy(eslice[offset:], id[:])
		offset += 16
	}
	return eslice[:], nil
}

func (e *Entities) UnmarshalBinary(data []byte) error {
	//TODO: this
	return nil
}

func (e Entities) MarshalText() (text []byte, err error) {
	slcs := make([]string, 0, len(e))
	for _, e := range e {
		slcs = append(slcs, encoderWorld.EntityUUID(e).String())
	}
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(slcs); err != nil {
		return nil, err
	}
	return []byte(buf.String()), nil
}

func (e *Entities) UnmarshalText(text []byte) error {
	slc := make([]string, 0)
	if _, err := toml.NewDecoder(bytes.NewReader(text)).Decode(&slc); err != nil {
		return err
	}
	eslc := make([]Entity, 0, len(slc))
	for _, s := range slc {
		id, err := uuid.Parse(s)
		if err == nil {
			eslc = append(eslc, decoderWorld.getEntityByUUID(id))
		}
	}
	v := Entities(eslc)
	*e = v
	return nil
}

func (e *Entities) UnmarshalJSON(b []byte) error {
	slc := make([]string, 0)
	if err := json.Unmarshal(b, &slc); err != nil {
		return err
	}
	eslc := make([]Entity, 0, len(slc))
	for _, s := range slc {
		id, err := uuid.Parse(s)
		if err == nil {
			eslc = append(eslc, decoderWorld.getEntityByUUID(id))
		}
	}
	v := Entities(eslc)
	*e = v
	return nil
}

func (e Entities) MarshalJSON() ([]byte, error) {
	slc := make([]string, 0, len(e))
	for _, e := range e {
		slc = append(slc, encoderWorld.EntityUUID(e).String())
	}
	return json.Marshal(slc)
}
