package command

import (
	"testing"

	"github.com/AutoDruid/photon-parser/internal/hooks"
	v18 "github.com/AutoDruid/photon-parser/internal/parameters/v18"
	"github.com/AutoDruid/photon-parser/internal/types"
)

func TestAsyncCommandEmitFullChannelDoesNotAllocate(t *testing.T) {
	h := hooks.NewHooks[v18.Parameter]()
	h.OnCommandAsync(types.HookOptions{Size: 1})
	h.AsyncHooks.OnCommand <- types.Command[v18.Parameter]{}

	cmd := types.Command[v18.Parameter]{
		CommandHeader: types.CommandHeader{Type: types.SendReliableCommand},
		ReliablePayload: types.Reliable[v18.Parameter]{
			Parameters: []v18.Parameter{{Header: v18.Header{ID: 1, Type: v18.Int8Type}}},
		},
	}

	allocs := testing.AllocsPerRun(100, func() {
		emit(h, &cmd)
	})

	if allocs != 0 {
		t.Fatalf("emit() allocations on full async channel = %v, want 0", allocs)
	}
}
