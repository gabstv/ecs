package ecs

import "fmt"

type Error struct {
	Code   int
	Reason string
}

func (e Error) Error() string {
	return fmt.Sprintf("%d - %s", e.Code, e.Reason)
}

var (
	ErrSystemInvalid = Error{1, "System is invalid"}
)
