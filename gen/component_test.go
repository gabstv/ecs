package gen

import (
	"testing"

	"github.com/dave/jennifer/jen"
)

func TestComponentGeneration(t *testing.T) {
	f := jen.NewFile("xyz")
	f.PackageComment("Code generated by ecs https://github.com/gabstv/ecs; DO NOT EDIT.")
	Component(f, ComponentDef{
		StructName: "Position",
	})
	println(f.GoString())
}
