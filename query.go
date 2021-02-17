package ecs

type MatchFn func(f Flag, w World) bool

// SortedEntities implements sort.Interface
type SortedEntities []Entity

func (a SortedEntities) Len() int           { return len(a) }
func (a SortedEntities) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortedEntities) Less(i, j int) bool { return a[i] < a[j] }

type EntityFlag struct {
	Entity Entity
	Flag   Flag
}

type EntityFlagSlice []EntityFlag

func (a EntityFlagSlice) Len() int           { return len(a) }
func (a EntityFlagSlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a EntityFlagSlice) Less(i, j int) bool { return a[i].Entity < a[j].Entity }
