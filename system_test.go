package ecs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testSharedRegistryData struct {
	IDs      []int
	TimesRun int
}

func TestGlobalRegistryFlow(t *testing.T) {

	RegisterGlobalSystem(GlobalSystemInfo[BenchPos3]{
		ExecPriority: 100,
		ExecFlag:     80,
		ExecBuilder: func(w *World, s *System[BenchPos3]) func(view *View[BenchPos3]) {
			data := s.Data().(*testSharedRegistryData)
			return func(view *View[BenchPos3]) {
				data.TimesRun++
			}
		},
		Initializer: func(w *World, sys *System[BenchPos3]) {
			data := &testSharedRegistryData{
				IDs: make([]int, 0, 2),
			}
			data.IDs = append(data.IDs, sys.ID())
			w.Data().Set("tdata", data)
			sys.SetData(data)
		},
	})
	RegisterGlobalSystem(GlobalSystemInfo[BenchPos3]{
		ExecPriority: 99,
		ExecFlag:     80,
		ExecBuilder: func(w *World, s *System[BenchPos3]) func(view *View[BenchPos3]) {
			data := s.Data().(*testSharedRegistryData)
			return func(view *View[BenchPos3]) {
				data.TimesRun++
			}
		},
		Initializer: func(w *World, sys *System[BenchPos3]) {
			data := w.Data().Get("tdata").(*testSharedRegistryData)
			data.IDs = append(data.IDs, sys.ID())
			sys.SetData(data)
		},
	})

	ww := NewWorld()
	xdata := ww.Data().Get("tdata").(*testSharedRegistryData)
	assert.Equal(t, 2, len(xdata.IDs))
	ww.StepF(80)
	assert.Equal(t, 2, xdata.TimesRun)
	ww.StepF(80)
	assert.Equal(t, 4, xdata.TimesRun)
}
