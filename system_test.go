package ecs

// import (
// 	"reflect"
// 	"testing"
// )

// func TestComparisons(t *testing.T) {
// 	cmd2 := &Commands{
// 		list: []Command{
// 			func() {},
// 		},
// 	}
// 	if reflect.TypeOf(cmd2) != typeCommandsPointer {
// 		t.Error("Commands pointer type not detected")
// 	}
// 	if reflect.TypeOf(cmd2).Elem() != typeCommandsPointer.Elem() {
// 		t.Error("Commands pointer type not detected")
// 	}
// 	type basics struct{}
// 	cmd4 := reflect.TypeOf(&basics{})
// 	if cmd4 == typeCommandsPointer {
// 		t.Error("Commands pointer type not detected")
// 	}
// }
