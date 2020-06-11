package ecs_test

import (
	"testing"
	"unsafe"

	"github.com/gabstv/ecs"
)

type TestMarco struct {
	X float64 // 8 bytes
	Y float64 // 8 bytes
	Z float64 // 8 bytes
}

type TestPolo struct {
	Scale float64 // 8 bytes
}

func BenchmarkViews100000(b *testing.B) {
	b.StopTimer()
	w := ecs.NewWorld()
	c1, _ := w.NewComponent(ecs.NewComponentInput{
		Name: "cmarco",
	})
	c2, _ := w.NewComponent(ecs.NewComponentInput{
		Name: "cpolo",
	})
	w.NewSystem("smarco", 100, func(ctx ecs.Context) {
		matches := ctx.System().View().Matches()
		for _, v := range matches {
			a := v.Components[c1].(*TestMarco)
			a.X += 0.1
			a.Y += 0.2
			a.Z += 0.3
		}
	}, c1)
	fmax := func(a, b, c float64) float64 {
		if a > b && a > c {
			return a
		}
		if b > c {
			return b
		}
		return c
	}
	w.NewSystem("smarcopolo", 99, func(ctx ecs.Context) {
		matches := ctx.System().View().Matches()
		for _, v := range matches {
			a := v.Components[c1].(*TestMarco)
			b := v.Components[c2].(*TestPolo)
			b.Scale = fmax(a.X, a.Y, a.Z)
			v.Components[c2] = b
		}
	}, c1, c2)
	entities := w.NewEntities(100000)
	for i, e := range entities {
		w.AddComponentToEntity(e, c1, &TestMarco{})
		if i%7 == 0 {
			w.AddComponentToEntity(e, c2, &TestPolo{})
		}
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		w.Run(1 / 60)
	}
}

const szmarco = int(unsafe.Sizeof(TestMarco{}))
const szpolo = int(unsafe.Sizeof(TestPolo{}))

func BenchmarkUnsafeParallelViews100000(b *testing.B) {
	b.StopTimer()
	//
	fmax := func(a, b, c float64) float64 {
		if a > b && a > c {
			return a
		}
		if b > c {
			return b
		}
		return c
	}
	//
	marcomap := make(map[ecs.Entity]int)
	polomap := make(map[ecs.Entity]int)
	marcodata := make([]byte, szmarco*100000)
	polodata := make([]byte, szpolo*(100000))
	//
	marcobytes := func(m TestMarco) []byte {
		b := *(*[szmarco]byte)(unsafe.Pointer(&m))
		return b[:]
	}
	polobytes := func(m TestPolo) []byte {
		b := *(*[szpolo]byte)(unsafe.Pointer(&m))
		return b[:]
	}
	//
	moff := 0
	poff := 0
	firstmap := make([]ecs.Entity, 0, (100000))
	combomap := make([]ecs.Entity, 0, (100000/7)+1)
	for i := 0; i < 100000; i++ {
		e := ecs.Entity(i + 1)
		copy(marcodata[moff:moff+szmarco], marcobytes(TestMarco{}))
		marcomap[e] = moff
		moff += szmarco
		firstmap = append(firstmap, e)
		if i%7 == 0 {
			copy(polodata[poff:poff+szpolo], polobytes(TestPolo{}))
			polomap[e] = poff
			poff += szpolo
			combomap = append(combomap, e)
		}
	}
	//
	mbuffer := [szmarco]byte{}
	marcob2s := func(v []byte, offset int) TestMarco {
		copy(mbuffer[:], v[offset:offset+szmarco])
		return *(*TestMarco)(unsafe.Pointer(&mbuffer))
	}
	//
	pbuffer := [szpolo]byte{}
	polob2s := func(v []byte, offset int) TestPolo {
		copy(mbuffer[:], v[offset:offset+szmarco])
		return *(*TestPolo)(unsafe.Pointer(&pbuffer))
	}
	//
	b.StartTimer()
	//
	for bi := 0; bi < b.N; bi++ {
		for i := 0; i < len(firstmap); i++ {
			e := firstmap[i]
			marco := marcob2s(marcodata, marcomap[e])
			marco.X += .1
			marco.Y += .2
			marco.Z += .3
			copy(marcodata[marcomap[e]:marcomap[e]+szmarco], marcobytes(marco))
		}
		for i := 0; i < len(combomap); i++ {
			e := combomap[i]
			marco := marcob2s(marcodata, marcomap[e])
			polo := polob2s(polodata, polomap[e])
			polo.Scale = fmax(marco.X, marco.Y, marco.Z)
			copy(polodata[polomap[e]:polomap[e]+szpolo], polobytes(polo))
		}
	}
}
