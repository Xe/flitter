package mock

import "testing"

func TestMakeBackend(t *testing.T) {
	b := &Backend{}

	if b == nil {
		t.Fatal("Backend is nil")
	}
}
