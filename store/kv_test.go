package store

import "testing"

func TestCommandSet_Apply(t *testing.T) {
	actor := NewActor()

	cmd := CommandSet{Key: "hello", Value: "world", Error: make(chan error)}
	err := cmd.Apply(actor)
	if err != nil {
		t.Errorf("Error applying command: %v", err)
	}

	// happy path
	if actor.store["hello"] != "world" {
		t.Errorf("Expected 'world' but got %s", actor.store["hello"])
	}

	// error case
	cmd = CommandSet{Key: "hello", Value: "world", Error: make(chan error)}
	err = cmd.Apply(actor)
	if err == nil {
		t.Errorf("Expected error but got nil")
	}
}

func TestCommandGet_Apply(t *testing.T) {
	actor := NewActor()

	cmdSet := CommandSet{Key: "hello", Value: "world", Error: make(chan error)}
	err := cmdSet.Apply(actor)
	if err != nil {
		t.Errorf("Error applying command: %v", err)
	}

	// happy path
	cmdGet := CommandGet{Key: "hello", Response: make(chan string, 1), Error: make(chan error)}
	go func() {
		err = cmdGet.Apply(actor)
		if err != nil {
			t.Errorf("Error applying command: %v", err)
		}
	}()

	select {
	case res := <-cmdGet.Response:
		if res != "world" {
			t.Errorf("Expected 'world' but got %s", res)
		}
	case err = <-cmdGet.Error:
		t.Errorf("Error getting key: %v", err)
	}

	// error case
	cmdGet = CommandGet{Key: "hello_notfound", Response: make(chan string, 1), Error: make(chan error)}
	go func() {
		err = cmdGet.Apply(actor)
		if err != nil {
			t.Errorf("Error applying command: %v", err)
		}
	}()

	select {
	case res := <-cmdGet.Response:
		t.Errorf("Expected empty response but got %s", res)
	case err = <-cmdGet.Error:
		if err == nil {
			t.Errorf("Expected error but got nil")
		}
	}
}

// TODO Performance/Benchmark tests
