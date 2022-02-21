package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnts(t *testing.T) {
	assert.Equal(t, 0, ent3(1, 1, 1))
	assert.Equal(t, 0, ent4(1, 1, 1, 1))
	assert.Equal(t, 1, ent3(1, 2, 3))
	assert.Equal(t, 1, ent4(1, 2, 3, 4))
	assert.Equal(t, 2, ent3(2, 1, 3))
	assert.Equal(t, 2, ent4(2, 1, 3, 4))
	assert.Equal(t, 3, ent3(3, 2, 1))
	assert.Equal(t, 3, ent4(3, 2, 1, 4))
	assert.Equal(t, 4, ent4(4, 3, 2, 1))
}

func TestInsert(t *testing.T) {
	slc1 := make([]int, 4, 8)
	slc1[0], slc1[1] = 1, 2
	slc1[2], slc1[3] = 3, 4
	slc1 = Insert(slc1, 0, 5, 6)
	assert.Equal(t, []int{5, 6, 1, 2, 3, 4}, slc1)
	slc1 = Insert(slc1, -1, 10, 20, 30)
	assert.Equal(t, []int{5, 6, 1, 2, 3, 4, 10, 20, 30}, slc1)
}
