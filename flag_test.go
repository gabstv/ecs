package ecs

import (
	"math/rand"
	"testing"
	"time"
)

// helper functions

func randstringarray(count int) []string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ll := make([]string, count)
	for i := 0; i < count; i++ {
		wlen := 3 + r.Intn(15)
		www := make([]byte, wlen)
		for k := range www {
			www[k] = byte(65 + r.Intn(25))
		}
		ll[i] = string(www)
	}
	return ll
}

func randflags(count int) []Flag {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ll := make([]Flag, count)
	for i := 0; i < count; i++ {
		wlen := uint8(r.Intn(256))
		ll[i] = NewFlag(wlen)
	}
	return ll
}

func randindexes(maxval, count int) []int {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ll := make([]int, count)
	for i := 0; i < count; i++ {
		ll[i] = r.Intn(maxval)
	}
	return ll
}

func TestFlagEquals(t *testing.T) {
	a := NewFlagRaw(1, 0, 1, 0)
	b := NewFlagRaw(1, 0, 0, 0)
	if a.Equals(b) {
		t.Fail()
	}
	c := NewFlagRaw(1, 0, 1, 0)
	if !a.Equals(c) {
		t.Fail()
	}
}

func TestFlagBitmap(t *testing.T) {
	bmap := NewFlag(1).Or(NewFlag(3)).Or(NewFlag(200))
	if bmap.Contains(NewFlag(2)) {
		t.Fail()
	}
	if bmap.Contains(NewFlag(199)) {
		t.Fail()
	}
	if !bmap.Contains(NewFlag(200)) {
		t.Fail()
	}
	if !bmap.Contains(NewFlag(1)) {
		t.Fail()
	}
}

func TestFlagContainsAny(t *testing.T) {
	big := NewFlag(1).Or(NewFlag(2)).Or(NewFlag(3)).Or(NewFlag(4)).Or(NewFlag(5))
	if !big.ContainsAny(NewFlag(4)) {
		t.Fail()
	}
	if big.ContainsAny(NewFlag(6)) {
		t.Fail()
	}
	if !big.ContainsAny(NewFlag(6).Or(NewFlag(2))) {
		t.Fail()
	}
	if !big.ContainsAny(NewFlag(4).Or(NewFlag(2))) {
		t.Fail()
	}
}

func BenchmarkStringMap128(b *testing.B) {
	b.StopTimer()
	set0 := randstringarray(128)
	indexes := randindexes(128, 512)
	mmap := make(map[string]bool)
	for _, v := range set0 {
		mmap[v] = true
	}
	b.StartTimer()
	for i := 0; i < b.N; i += 512 {
		for j := 0; j < 512; j++ {
			_ = mmap[set0[indexes[j]]]
		}
	}
}

func BenchmarkFlag128(b *testing.B) {
	b.StopTimer()
	set0 := randflags(128)
	indexes := randindexes(128, 512)
	master := NewFlagRaw(0, 0, 0, 0)
	for _, v := range set0 {
		master = master.Or(v)
	}
	b.StartTimer()
	for i := 0; i < b.N; i += 512 {
		for j := 0; j < 512; j++ {
			_ = master.Contains(set0[indexes[j]])
		}
	}
}
