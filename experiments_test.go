package ecs

import (
	"math"
	"math/rand"
	"strconv"
	"testing"
)

func BenchmarkComponentSliceCopy(b *testing.B) {
	b.StopTimer()
	type Object struct {
		Name  string
		X     float64
		Y     float64
		Rot   float64
		Score int
	}
	objects := make([]Object, 1000000)
	objcopy := make([]Object, 1000000)
	for i := 0; i < 1000000; i++ {
		objects[i] = Object{
			Name:  "object_" + strconv.Itoa(i),
			X:     rand.Float64() * 100.0,
			Y:     rand.Float64() * 100.0,
			Rot:   rand.Float64() * math.Pi * 2.0,
			Score: rand.Int(),
		}
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		copy(objcopy, objects)
	}
}
