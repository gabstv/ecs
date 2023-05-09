package ecs

import (
	// "bytes"
	//"hash/crc32"
	"reflect"
)

// type TypeHash uint32

// func newTypeHash(types ...reflect.Type) TypeHash {
// 	tape := new(bytes.Buffer)
// 	for _, t := range types {
// 		tape.WriteString(t.PkgPath())
// 		tape.WriteString(t.Name())
// 	}
// 	return TypeHash(crc32.ChecksumIEEE(tape.Bytes()))
// }

func zeroValue(t reflect.Type) reflect.Value {
	return reflect.Zero(t)
}

// TypeTape is a map key of a query.
type TypeTape [16]reflect.Value

type TypeMapKey reflect.Value

func typeMapKeyOf(t reflect.Type) TypeMapKey {
	return TypeMapKey(zeroValue(t))
}

func typeTapeOf(types ...reflect.Type) TypeTape {
	if len(types) > 16 {
		panic("too many types (limit = 16)")
	}
	tape := TypeTape{}
	for i, t := range types {
		tape[i] = zeroValue(t)
	}
	return tape
}
