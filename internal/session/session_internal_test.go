package session

import (
	"testing"

	"github.com/AutoDruid/photon-parser/internal/hooks"
	v18 "github.com/AutoDruid/photon-parser/internal/parameters/v18"
	"github.com/AutoDruid/photon-parser/internal/types"
)

func TestAsyncSessionEmitFullChannelDoesNotAllocate(t *testing.T) {
	h := hooks.NewHooks[v18.Parameter]()
	h.OnSessionAsync(types.HookOptions{Size: 1})
	h.AsyncHooks.OnSession <- types.Session[v18.Parameter]{}

	sess := types.Session[v18.Parameter]{
		Commands: []types.Command[v18.Parameter]{
			{
				CommandHeader: types.CommandHeader{Type: types.SendReliableCommand},
				ReliablePayload: types.Reliable[v18.Parameter]{
					Parameters: []v18.Parameter{{Header: v18.Header{ID: 1, Type: v18.Int8Type}}},
				},
			},
		},
	}

	allocs := testing.AllocsPerRun(100, func() {
		emit(h, &sess)
	})

	if allocs != 0 {
		t.Fatalf("emit() allocations on full async channel = %v, want 0", allocs)
	}
}
