package store

import (
	"bytes"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

// Command tests

var wg sync.WaitGroup

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

func TestActor_ConcurrentSet(t *testing.T) {
	// arrange
	actor := NewActor()

	// act
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			actor.Set("key"+strconv.Itoa(i), "world")
			wg.Done()
		}()
	}
	wg.Wait()

	// assert
	list := actor.List()
	if len(list) != 10 {
		t.Errorf("Expected 10 keys but got %v", list)
	}
}

func TestActor_ConcurrentSetSameKey(t *testing.T) {
	// arrange
	actor := NewActor()
	actor.quiet = true
	key := "hello"

	// act
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			actor.Set(key, "world")
			wg.Done()
		}()
	}
	wg.Wait()

	// assert
	val, err := actor.Get(key)
	if err != nil || val != "world" {
		t.Errorf("Unexpected error or value after Set: %v, %s", err, val)
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

func TestActor_GetError(t *testing.T) {
	// arrange
	actor := NewActor()

	// act
	val, err := actor.Get("hello")

	// assert
	if err == nil || val != "" {
		t.Errorf("Expected error but got nil or value: %v, %s", err, val)
	}
}

func TestActor_ConcurrentGet(t *testing.T) {
	// arrange
	actor := NewActor()
	actor.store["hello"] = "world"
	actor.quiet = true
	key := "hello"

	// act
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			actor.Get(key)
			wg.Done()
		}()
	}
	wg.Wait()

	// assert
	val, err := actor.Get(key)
	if err != nil || val != "world" {
		t.Errorf("Unexpected error or value after Set: %v, %s", err, val)
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

func TestActor_ConcurrentDelete(t *testing.T) {
	// arrange
	actor := NewActor()
	actor.store["hello"] = "world"
	key := "hello"

	// act
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			actor.Delete(key)
			wg.Done()
		}()
	}
	wg.Wait()

	// assert
	if _, ok := actor.store[key]; ok {
		t.Errorf("Expected key to be deleted but still exists")
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

func TestNewActor(t *testing.T) {
	// act
	actor := NewActor()

	// assert
	if actor == nil {
		t.Errorf("Expected actor to be created but got nil")
	}

	if actor.store == nil {
		t.Errorf("Expected store to be created but got nil")
	}

}

// Benchmarks

func BenchmarkActor_Set(b *testing.B) {
	actor := NewActor()
	actor.quiet = true

	for i := 0; i < b.N; i++ {
		actor.Set("key"+strconv.Itoa(i), "value")
	}
}

func BenchmarkActor_SetConcurrent(b *testing.B) {
	actor := NewActor()
	actor.quiet = true

	b.RunParallel(func(pb *testing.PB) {
		for i := 0; pb.Next(); i++ {
			actor.Set("key"+strconv.Itoa(i), "value")
		}
	})
}

func BenchmarkActor_SetConcurrentSameKey(b *testing.B) {
	actor := NewActor()
	actor.quiet = true

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			actor.Set("key", "value")
		}
	})
}

func BenchmarkActor_GetConcurrent(b *testing.B) {
	actor := NewActor()

	for i := 0; i < b.N; i++ {
		actor.Set("key"+strconv.Itoa(i), "value")
	}

	b.RunParallel(func(pb *testing.PB) {
		for i := 0; pb.Next(); i++ {
			actor.Get("key" + strconv.Itoa(i))
		}
	})
}

func BenchmarkActor_GetConcurrentSameKey(b *testing.B) {
	actor := NewActor()
	actor.store["key"] = "value"
	actor.quiet = true

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			actor.Get("key")
		}
	})
}

func BenchmarkActor_DeleteConcurrent(b *testing.B) {
	actor := NewActor()
	actor.quiet = true

	for i := 0; i < b.N; i++ {
		actor.Set("key"+strconv.Itoa(i), "value")
	}

	b.RunParallel(func(pb *testing.PB) {
		for i := 0; pb.Next(); i++ {
			actor.Delete("key" + strconv.Itoa(i))
		}
	})
}

func BenchmarkActor_DeleteConcurrentSameKey(b *testing.B) {
	actor := NewActor()
	actor.store["key"] = "value"
	actor.quiet = true

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			actor.Delete("key")
		}
	})
}

func BenchmarkActor_ConcurrentMixedOperations(b *testing.B) {
	actor := NewActor()
	actor.quiet = true

	for i := 0; i < b.N; i++ {
		actor.Set("key"+strconv.Itoa(i), "value")
	}

	b.RunParallel(func(pb *testing.PB) {
		for i := 0; pb.Next(); i++ {
			switch i % 3 {
			case 0:
				actor.Set("key"+strconv.Itoa(i), "value")
			case 1:
				actor.Get("key" + strconv.Itoa(i))
			case 2:
				actor.Delete("key" + strconv.Itoa(i))
			}
		}
	})
}

func BenchmarkActor_StressTestSet(b *testing.B) {
	actor := NewActor()
	actor.quiet = true

	b.SetBytes(10000)
	b.RunParallel(func(pb *testing.PB) {
		for i := 0; pb.Next(); i++ {
			actor.Set("key"+strconv.Itoa(i), strings.Repeat("value", 1000))
		}
	})
}

// Notes from tests/benchmarks:
// after trying both GWT and AAA styles, I found that AAA was more readable and easier to follow for me
// (also thinking from perspective of someone else reading the code)
//
// use a bool flag to suppress logs in tests, helpful to not flood console
// you can use sync waitgroups in tests when testing concurrent operations
// - helpful when you need to wait for all routines to finish
//
// buffers can be used to capture logs in tests we expect during error conditions
// - helped with asserting the expected error messages
//
// the b.RunParallel func is used when you want to run multiple go routines in parallel
// - this is great for simulating a proper real world scenario where operations can come in concurrently
//
// when testing a mixture of operations, using a switch statement helped to randomly select an operation
// - another great way to simulate real world scenarios where all sorts of ops are coming in
//
// the b.SetBytes method is great for setting the "size" of the operation being tested
// - simulating the program under high loads by having it process large amounts of data
// - the repeated string is a way to simulate larger amounts of data (putting more stress on the actor.Set)
