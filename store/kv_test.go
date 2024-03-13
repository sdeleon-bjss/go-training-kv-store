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
	// given
	actor := NewActor()

	// when
	cmd := CommandSet{Key: "hello", Value: "world", Error: make(chan error)}
	err := cmd.Apply(actor)
	if err != nil {
		t.Errorf("Error applying command: %v", err)
	}

	// then
	if actor.store["hello"] != "world" {
		t.Errorf("Expected 'world' but got %s", actor.store["hello"])
	}

	// given error
	cmd = CommandSet{Key: "hello", Value: "world", Error: make(chan error)}
	err = cmd.Apply(actor)

	// then
	if err == nil {
		t.Errorf("Expected error but got nil")
	}
}

func TestCommandGet_Apply(t *testing.T) {
	// given
	actor := NewActor()

	cmdSet := CommandSet{Key: "hello", Value: "world", Error: make(chan error)}
	err := cmdSet.Apply(actor)
	if err != nil {
		t.Errorf("Error applying command: %v", err)
	}

	// when
	cmdGet := CommandGet{Key: "hello", Response: make(chan string, 1), Error: make(chan error)}
	go func() {
		err = cmdGet.Apply(actor)
		if err != nil {
			t.Errorf("Error applying command: %v", err)
		}
	}()

	// then
	select {
	case res := <-cmdGet.Response:
		if res != "world" {
			t.Errorf("Expected 'world' but got %s", res)
		}
	case err = <-cmdGet.Error:
		t.Errorf("Error getting key: %v", err)
	}

	// given error
	cmdGet = CommandGet{Key: "hello_notfound", Response: make(chan string, 1), Error: make(chan error)}
	go func() {
		err = cmdGet.Apply(actor)
		if err != nil {
			t.Errorf("Error applying command: %v", err)
		}
	}()

	// then
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
	// given
	actor := NewActor()

	cmdSet := CommandSet{Key: "hello", Value: "world", Error: make(chan error)}
	err := cmdSet.Apply(actor)
	if err != nil {
		t.Errorf("Error applying command: %v", err)
	}

	// when
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

func TestActor_Get(t *testing.T) {
	actor := NewActor()

	go actor.Set("hello", "world")
	time.Sleep(time.Second)

	val, err := actor.Get("hello")
	if err != nil || val != "world" {
		t.Errorf("Unexpected error or value after Set: %v, %s", err, val)
	}

	// error case
	val, err = actor.Get("hello_notfound")
	if err == nil || val != "" {
		t.Errorf("Expected error but got nil or value: %v, %s", err, val)
	}
}

func TestActor_Delete(t *testing.T) {
	// buffer
	buf := &bytes.Buffer{}
	log.SetOutput(buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	actor := NewActor()

	go actor.Set("hello", "world")
	time.Sleep(time.Second)

	go actor.Delete("hello")
	time.Sleep(time.Second)

	val, err := actor.Get("hello")
	if err == nil || val != "" {
		t.Errorf("Expected error but got nil or value: %v, %s", err, val)
	}

	// error case
	go actor.Delete("hello_notfound")
	time.Sleep(time.Second)

	logs := buf.String()
	if !strings.Contains(logs, "Error applying command: key hello_notfound does not exist") {
		t.Errorf("Expected an error log for key %q but got: %s", "hello_notfound", logs)
	}
}

func TestActor_List(t *testing.T) {
	actor := NewActor()

	go actor.Set("hello", "world")
	go actor.Set("hello2", "world2")
	time.Sleep(time.Second)

	list := actor.List()
	if len(list) != 2 {
		t.Errorf("Expected 2 keys but got %v", list)
	}

	if list["hello"] != "world" || list["hello2"] != "world2" {
		t.Errorf("Unexpected key-value pairs: %v", list)
	}
}

// TODO Performance/Benchmark tests
