package photon

import "testing"

func TestDefaultConfigDoesNotAllocate(t *testing.T) {
	allocs := testing.AllocsPerRun(100, func() {
		_ = defaultConfig()
	})

	if allocs != 0 {
		t.Fatalf("defaultConfig() allocations = %v, want 0", allocs)
	}
}
