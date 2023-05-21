package ecs

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"testing"
	"unsafe"
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

func TestRealloc(t *testing.T) {
	slcs := make([][]int, 6)
	slcs[0] = make([]int, 2, 4)
	slcs[1] = make([]int, 1, 4)
	slcs[2] = make([]int, 3, 4)
	slcs[3] = make([]int, 0, 4)
	slcs[4] = make([]int, 1, 4)
	slcs[5] = make([]int, 1, 4)
	for i := 0; i < 6; i++ {
		for j := 0; j < 4; j++ {
			if len(slcs[i]) > j {
				slcs[i][j] = (1 + i) * (1 + j)
			}
		}
	}
	p00 := uintptr(unsafe.Pointer(&slcs[0][0]))
	p01 := uintptr(unsafe.Pointer(&slcs[5][0]))
	fmt.Println("p00:", p00, "cap", cap(slcs[0]))
	fmt.Println("p01:", p01, "cap", cap(slcs[5]))
	fmt.Println("dist:", p01-p00)
	fmt.Println("s0", unsafe.Sizeof(slcs[0]))
	fmt.Println("s1", unsafe.Sizeof(slcs[5]))
	slcp0 := slcs[0]
	slcp1 := slcs[5]

	slcs[0] = make([]int, len(slcp0), cap(slcp0)*2)
	slcs[5] = make([]int, len(slcp1), cap(slcp1)*2)
	copy(slcs[0], slcp0)
	copy(slcs[5], slcp1)
	slcp0 = nil
	slcp1 = nil
	p00 = uintptr(unsafe.Pointer(&slcs[0][0]))
	p01 = uintptr(unsafe.Pointer(&slcs[5][0]))
	fmt.Println("p00:", p00, "cap", cap(slcs[0]))
	fmt.Println("p01:", p01, "cap", cap(slcs[5]))
	fmt.Println("dist:", p01-p00)
	fmt.Println("s0", unsafe.Sizeof(slcs[0]))
	fmt.Println("s1", unsafe.Sizeof(slcs[5]))
	fmt.Println(unsafe.Sizeof(make([]int, 1)))
	p000 := unsafe.Pointer(&slcs[0])
	pint := (*int)(p000)
	pint2 := (*int)(unsafe.Pointer(uintptr(p000) + unsafe.Sizeof(uintptr(0))))
	pint3 := (*int)(unsafe.Pointer(uintptr(p000) + unsafe.Sizeof(uintptr(0))*2))
	fmt.Println("address of [8]int (slcs[0])", *pint)
	fmt.Println("size of [8]int (slcs[0])", *pint2)
	fmt.Println("capacity of [8]int (slcs[0])", *pint3)
	fmt.Println(uintptr(p000))
	fmt.Println("OK")
}
