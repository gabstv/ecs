package gen

import (
	"strings"

	"github.com/dave/jennifer/jen"
)

func mutexLU(enabled bool, paths string, xdefer bool) *jen.Statement {
	if !enabled {
		return jen.Empty()
	}
	var x *jen.Statement
	if paths == "" {
		panic("paths required")
	}
	ps := strings.Split(paths, ".")
	if xdefer {
		x = jen.Defer().Id(ps[0])
	} else {
		x = jen.Id(ps[0])
	}
	for _, v := range ps[1:] {
		x = x.Op(v)
	}
	x = x.Call()
	return x
}
