package ecs

import (
	"github.com/Pilatuz/bigz/uint256"
	"golang.org/x/exp/slices"
)

type U256 = uint256.Uint256

type Entity uint64

type fatEntity struct {
	Entity       Entity
	ComponentMap U256
	IsRemoved    bool
}

type EntityList struct {
	items []entityListItem
}

type entityListItem struct {
	Entity    Entity
	IsRemoved bool
}

func eliBinSearch(eli entityListItem, entity Entity) int {
	if eli.Entity < entity {
		return -1
	}
	if eli.Entity > entity {
		return 1
	}
	return 0
}

func (e *EntityList) Len() int {
	return len(e.items)
}

func (e *EntityList) Add(entity Entity) {
	if len(e.items) < 1 || e.items[len(e.items)-1].Entity < entity {
		e.items = append(e.items, entityListItem{
			Entity:    entity,
			IsRemoved: false,
		})
		return
	}
	index, ok := slices.BinarySearchFunc(e.items, entity, eliBinSearch)
	if ok {
		e.items[index].IsRemoved = false
		return
	}
	e.items = slices.Insert(e.items, index, entityListItem{
		Entity:    entity,
		IsRemoved: false,
	})
}

func (e *EntityList) Remove(entity Entity) {
	index, ok := slices.BinarySearchFunc(e.items, entity, eliBinSearch)
	if !ok {
		return
	}
	e.items[index].IsRemoved = true
}

func (e *EntityList) GC() {
	var newItems []entityListItem
	for _, v := range e.items {
		if !v.IsRemoved {
			newItems = append(newItems, v)
		}
	}
	e.items = newItems
}

func (e *EntityList) All() []Entity {
	el := make([]Entity, 0, len(e.items))
	for _, v := range e.items {
		if !v.IsRemoved {
			el = append(el, v.Entity)
		}
	}
	return el
}
