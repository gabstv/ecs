package gentests

import (
	"github.com/gabstv/ecs/v3"
)

type Position struct {
	X float64
	Y float64
}

//go:generate go run ../cmd/ecsgen/main.go -n Position -p gentests -o position_component.go --component-tpl --vars "UUID=3DF7F486-807D-4CE8-A187-37CED338137B"

type Rotation struct {
	Angle float64
}

//go:generate go run ../cmd/ecsgen/main.go -n Rotation -p gentests -o rotation_component.go --component-tpl --vars "UUID=56890133-3769-477A-B163-412C5ECC6B07" --vars "Cap=2"

//go:generate go run ../cmd/ecsgen/main.go -n PosRot -p gentests -o posrot_system.go --system-tpl --vars "Priority=0" --vars "UUID=58FFC3BE-7BC8-4381-A93B-74945405F171" --components "Position" --components "Rotation"

var matchPosRotSystem = func(eflag ecs.Flag, w ecs.BaseWorld) bool {
	return eflag.Contains(GetPositionComponent(w).flag.Or(GetRotationComponent(w).Flag()))
}

var resizematchPosRotSystem = func(eflag ecs.Flag, w ecs.BaseWorld) bool {
	if eflag.Contains(GetPositionComponent(w).Flag()) {
		return true
	}
	if eflag.Contains(GetRotationComponent(w).Flag()) {
		return true
	}
	return false
}
