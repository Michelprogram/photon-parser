package context

import "testing"

func TestPoolGetAllocatesToRequestedSizeOnFirstUse(t *testing.T) {
	t.Parallel()

	pool := NewPool[int](64)
	items := pool.Get(1)

	if cap(items.Items) != 1 {
		t.Fatalf("cap(items.Items) = %d, want 1", cap(items.Items))
	}
}
