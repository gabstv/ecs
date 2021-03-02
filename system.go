package ecs

type System interface {
	UUID() string
	Name() string
	ComponentAdded(e Entity, eflag Flag)
	ComponentRemoved(e Entity, eflag Flag)
	ComponentWillResize(cflag Flag)
	ComponentResized(cflag Flag)
	//V() View
	Priority() int64
	Setup(w World)
	Enable()
	Disable()
	Enabled() bool
}

// BaseSystem implements Enable(), Disable() and Enabled() from ecs.System
type BaseSystem struct {
	enabled bool
}

// Enable this system
func (s *BaseSystem) Enable() {
	s.enabled = true
}

// Disable this system
func (s *BaseSystem) Disable() {
	s.enabled = false
}

// Enabled returns if this system is enabled
func (s *BaseSystem) Enabled() bool {
	return s.enabled
}
