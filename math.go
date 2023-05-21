package ecs

import "golang.org/x/exp/constraints"

func max[T constraints.Integer](a, b T) T {
	if a > b {
		return a
	}
	return b
}
