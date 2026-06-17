package trick

import "testing"

func TestNewMyTicker(t *testing.T) {
	called := false
	tick := NewMyTick(1, func() {
		called = true
	})
	defer tick.MyTick.Stop()

	if tick.MyTick == nil {
		t.Fatal("expected ticker to be initialized")
	}
	tick.Runner()
	if !called {
		t.Fatal("expected runner to be callable")
	}
}
