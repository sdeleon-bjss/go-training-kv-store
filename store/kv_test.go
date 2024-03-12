package store

import "testing"

func TestCommandSet_Apply(t *testing.T) {
	actor := NewActor()
	cmd := CommandSet{Key: "hello", Value: "world", Error: make(chan error)}

	err := cmd.Apply(actor)
	if err != nil {
		t.Errorf("Error applying command: %v", err)
	}

	if actor.store["hello"] != "world" {
		t.Errorf("Expected 'world' but got %s", actor.store["hello"])
	}
}

// TODO Performance/Benchmark tests
