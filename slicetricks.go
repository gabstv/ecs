package ecs

import "sort"

// Insert only copies elements
// in a[i:] once and allocates at most once.
// But, as of Go toolchain 1.16, due to lacking of
// optimizations to avoid elements clearing in the
// "make" call, the verbose way is not always faster.
//
// Future compiler optimizations might implement
// both in the most efficient ways.
func Insert[T any](s []T, k int, vs ...T) []T {
	if k >= len(s) || k == -1 {
		return append(s, vs...)
	}
	if n := len(s) + len(vs); n <= cap(s) {
		s2 := s[:n]
		copy(s2[k+len(vs):], s[k:])
		copy(s2[k:], vs)
		return s2
	}
	s2 := make([]T, len(s)+len(vs))
	copy(s2, s[:k])
	copy(s2[k:], vs)
	copy(s2[k+len(vs):], s[k:])
	return s2
}

func ent3(e1, e2, e3 Entity) int {
	if e1 == e2 && e1 == e3 {
		return 0
	}
	if e1 < e2 {
		if e1 < e3 {
			return 1
		}
		return 3
	}
	if e2 < e3 {
		return 2
	}
	return 3
}

// ent4 returns the index of the smallest entity.
func ent4(e1, e2, e3, e4 Entity) int {
	if e1 == e2 && e1 == e3 && e1 == e4 {
		return 0
	}
	if e1 < e2 {
		if e1 < e3 {
			if e1 < e4 {
				return 1
			}
			return 4
		}
		if e3 < e4 {
			return 3
		}
		return 4
	}
	if e2 < e3 {
		if e2 < e4 {
			return 2
		}
		return 4
	}
	if e3 < e4 {
		return 3
	}
	return 4
}

type Sortable[TI comparable, TD any] struct {
	Index TI
	Data  TD
}

// RemoveEntityFromSlice removes an entity from a slice. This func is useful
// when you don't know if the slice is sorted.
func RemoveEntityFromSlice(cls []Entity, e Entity) []Entity {
	for i, c := range cls {
		if c == e {
			return append(cls[:i], cls[i+1:]...)
		}
	}
	return cls
}

// AddEntityUnique adds an entity to a slice if it is not already in the slice.
func AddEntityUnique(cls []Entity, e Entity) []Entity {
	for _, c := range cls {
		if c == e {
			return cls
		}
	}
	return append(cls, e)
}

func SortEntities(s []Entity) {
	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j]
	})
}
