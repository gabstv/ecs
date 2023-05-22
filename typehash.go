package ecs

// "bytes"
//"hash/crc32"

// TypeTape is a map key of a query.
type TypeTape [16]ComponentUUID

func typeTapeOf(types ...Component) TypeTape {
	if len(types) > 16 {
		panic("too many types (limit = 16)")
	}
	tape := TypeTape{}
	for i, t := range types {
		tape[i] = t.ComponentUUID()
	}
	return tape
}
