# Entity Component System

A fast, code generate ECS (no more interface{}). Game Engine agnostic.

[![GoDoc](https://godoc.org/github.com/gabstv/ecs?status.svg)](https://godoc.org/github.com/gabstv/ecs)

```
// Auto-generated variant of sort.go:insertionSort
    10  func insertionSort_func(data lessSwap, a, b int) {
    11  	for i := a + 1; i < b; i++ {
    12  		for j := i; j > a && data.Less(j, j-1); j-- {
    13  			data.Swap(j, j-1)
    14  		}
    15  	}
    16  }
    ```