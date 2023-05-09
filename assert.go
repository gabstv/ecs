package ecs

func assert(t bool, msg string) {
	if !t {
		panic("assert: " + msg)
	}
}
