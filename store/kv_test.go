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
	// arrange
	actor := NewActor()

	// act
	cmd := CommandSet{Key: "hello", Value: "world", Error: make(chan error)}
	err := cmd.Apply(actor)
	if err != nil {
		t.Errorf("Error applying command: %v", err)
	}

	// assert
	if actor.store["hello"] != "world" {
		t.Errorf("Expected 'world' but got %s", actor.store["hello"])
	}
}

func TestCommandSet_ApplyError(t *testing.T) {
	// arrange
	actor := NewActor()
	actor.store["hello"] = "world"

	// act
	cmd := CommandSet{Key: "hello", Value: "world", Error: make(chan error)}
	err := cmd.Apply(actor)

	// assert
	if err == nil {
		t.Errorf("Expected error but got nil")
	}
}

func TestCommandGet_Apply(t *testing.T) {
	// arrange
	actor := NewActor()
	actor.store["hello"] = "world"

	// act
	cmdGet := CommandGet{Key: "hello", Response: make(chan string, 1), Error: make(chan error)}
	go func() {
		err := cmdGet.Apply(actor)
		if err != nil {
			t.Errorf("Error applying command: %v", err)
		}
	}()

	// assert
	select {
	case res := <-cmdGet.Response:
		if res != "world" {
			t.Errorf("Expected 'world' but got %s", res)
		}
	case err := <-cmdGet.Error:
		t.Errorf("Error getting key: %v", err)
	}
}

func TestCommandGet_ApplyError(t *testing.T) {
	// arrange
	actor := NewActor()
	actor.store["hello"] = "world"

	// act
	cmdGet := CommandGet{Key: "hello_notfound", Response: make(chan string, 1), Error: make(chan error)}
	go func() {
		err := cmdGet.Apply(actor)
		if err == nil {
			t.Errorf("Expected error but got nil")
		}
	}()

	// assert
	select {
	case res := <-cmdGet.Response:
		t.Errorf("Expected empty response but got %s", res)
	case err := <-cmdGet.Error:
		if err == nil {
			t.Errorf("Expected error but got nil")
		}
	}
}

func TestCommandDelete_Apply(t *testing.T) {
	// arrange
	actor := NewActor()
	actor.store["hello"] = "world"

	// act
	cmdDel := CommandDelete{Key: "hello", Error: make(chan error)}
	err := cmdDel.Apply(actor)
	if err != nil {
		t.Errorf("Error applying command: %v", err)
	}

	// assert
	if _, ok := actor.store["hello"]; ok {
		t.Errorf("Expected key to be deleted but still exists")
	}
}

func TestCommandDelete_ApplyError(t *testing.T) {
	// arrange
	actor := NewActor()

	// act
	cmdDel := CommandDelete{Key: "hello_notfound", Error: make(chan error)}
	err := cmdDel.Apply(actor)

	// assert
	if err == nil {
		t.Errorf("Expected error but got nil")
	}
}

// Actor tests

func TestActor_Set(t *testing.T) {
	// arrange
	actor := NewActor()

	// act
	go actor.Set("hello", "world")
	time.Sleep(time.Second)

	// assert
	val, err := actor.Get("hello")
	if err != nil || val != "world" {
		t.Errorf("Unexpected error or value after Set: %v, %s", err, val)
	}
}

func TestActor_SetError(t *testing.T) {
	// buffer
	buffer := &bytes.Buffer{}
	log.SetOutput(buffer)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	// arrange
	actor := NewActor()

	// act
	go actor.Set("hello", "world")
	time.Sleep(time.Second)

	// act
	go actor.Set("hello", "world again")
	time.Sleep(time.Second)

	// assert
	logs := buffer.String()
	if !strings.Contains(logs, "Error applying command: key hello already exists") {
		t.Errorf("Expected error message but got %s", logs)
	}
}

func TestActor_Get(t *testing.T) {
	// arrange
	actor := NewActor()
	actor.store["hello"] = "world"

	// act
	val, err := actor.Get("hello")

	// assert
	if err != nil || val != "world" {
		t.Errorf("Unexpected error or value after Set: %v, %s", err, val)
	}
}

func TestAction_GetError(t *testing.T) {
	// arrange
	actor := NewActor()

	// act
	val, err := actor.Get("hello")

	// assert
	if err == nil || val != "" {
		t.Errorf("Expected error but got nil or value: %v, %s", err, val)
	}
}

func TestActor_Delete(t *testing.T) {
	// arrange
	actor := NewActor()
	actor.store["hello"] = "world"

	// act
	go actor.Delete("hello")
	time.Sleep(time.Second)

	// assert
	if _, ok := actor.store["hello"]; ok {
		t.Errorf("Expected key to be deleted but still exists")
	}
}

func TestActor_DeleteError(t *testing.T) {
	// buffer
	buffer := &bytes.Buffer{}
	log.SetOutput(buffer)
	defer func() {
		log.SetOutput(os.Stderr)
	}()

	// arrange
	actor := NewActor()

	// act
	go actor.Delete("hello")
	time.Sleep(time.Second)

	// assert
	logs := buffer.String()
	if !strings.Contains(logs, "Error applying command: key hello does not exist") {
		t.Errorf("Expected error message but got %s", logs)
	}
}

func TestActor_List(t *testing.T) {
	// arrange
	actor := NewActor()
	actor.store["hello"] = "world"
	actor.store["hello2"] = "world2"

	// act
	list := actor.List()

	// assert
	if len(list) != 2 {
		t.Errorf("Expected 2 keys but got %v", list)
	}

	if list["hello"] != "world" || list["hello2"] != "world2" {
		t.Errorf("Unexpected key-value pairs: %v", list)
	}
}

// TODO Performance/Benchmark tests
