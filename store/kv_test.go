package store

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
	"time"
)

// Command tests

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

func TestCommandDelete_Apply(t *testing.T) {
	actor := NewActor()

	cmdSet := CommandSet{Key: "hello", Value: "world", Error: make(chan error)}
	err := cmdSet.Apply(actor)
	if err != nil {
		t.Errorf("Error applying command: %v", err)
	}

	// happy path
	cmdDel := CommandDelete{Key: "hello", Error: make(chan error)}
	err = cmdDel.Apply(actor)
	if err != nil {
		t.Errorf("Error applying command: %v", err)
	}

	// error case
	cmdDel = CommandDelete{Key: "hello_notfound", Error: make(chan error)}
	err = cmdDel.Apply(actor)
	if err == nil {
		t.Errorf("Expected error but got nil")
	}

	// check if key was deleted
	cmdGet := CommandGet{Key: "hello", Response: make(chan string, 1), Error: make(chan error)}
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

// Actor tests

func TestActor_Set(t *testing.T) {
	// buffer
	buf := &bytes.Buffer{}
	log.SetOutput(buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	actor := NewActor()

	// happy path
	go actor.Set("hello", "world")
	time.Sleep(time.Second)

	val, err := actor.Get("hello")
	if err != nil || val != "world" {
		t.Errorf("Unexpected error or value after Set: %v, %s", err, val)
	}

	// error case
	go actor.Set("hello", "world2")
	time.Sleep(time.Second)

	logs := buf.String()
	if !strings.Contains(logs, "Error applying command: key hello already exists") {
		t.Errorf("Expected error message but got %s", logs)
	}
}

// TODO Performance/Benchmark tests
