package ecs

import (
	"github.com/Pilatuz/bigz/uint256"
)

type U256 = uint256.Uint256

type Entity uint64

type fatEntity struct {
	Entity       Entity
	ComponentMap U256
	IsRemoved    bool
}
