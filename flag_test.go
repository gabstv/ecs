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

func randflags(count int) []flag {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ll := make([]flag, count)
	for i := 0; i < count; i++ {
		wlen := uint8(r.Intn(256))
		ll[i] = newflagbit(wlen)
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
	a := newflag(1, 0, 1, 0)
	b := newflag(1, 0, 0, 0)
	if a.equals(b) {
		t.Fail()
	}
	c := newflag(1, 0, 1, 0)
	if !a.equals(c) {
		t.Fail()
	}
}

func TestFlagBitmap(t *testing.T) {
	bmap := newflagbit(1).or(newflagbit(3)).or(newflagbit(200))
	if bmap.contains(newflagbit(2)) {
		t.Fail()
	}
	if bmap.contains(newflagbit(199)) {
		t.Fail()
	}
	if !bmap.contains(newflagbit(200)) {
		t.Fail()
	}
	if !bmap.contains(newflagbit(1)) {
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
	master := newflag(0, 0, 0, 0)
	for _, v := range set0 {
		master = master.or(v)
	}
	b.StartTimer()
	for i := 0; i < b.N; i += 512 {
		for j := 0; j < 512; j++ {
			_ = master.contains(set0[indexes[j]])
		}
	}
}
