package gentests

type Position struct {
	X float64
	Y float64
}

//go:generate go run ../cmd/ecsgen/main.go -n Position -p gentests -o position_component.go --component-tpl --vars "UUID=3DF7F486-807D-4CE8-A187-37CED338137B"

type Rotation struct {
	Angle float64
}
