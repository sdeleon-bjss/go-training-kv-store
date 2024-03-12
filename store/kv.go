package store

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Commands

type Command interface {
	Apply(*Actor) error
}

type CommandSet struct {
	Key   string
	Value string
	Error chan error
}

type CommandGet struct {
	Key      string
	Response chan string
	Error    chan error
}

type CommandDelete struct {
	Key   string
	Error chan error
}

func (c CommandSet) Apply(a *Actor) error {
	_, exists := a.store[c.Key]
	if exists {
		return fmt.Errorf("key %s already exists", c.Key)
	}
	//println("key val: ", c.Key, c.Value)

	a.store[c.Key] = c.Value
	return nil
}

func (c CommandGet) Apply(a *Actor) error {
	val, ok := a.store[c.Key]
	if !ok {
		c.Error <- fmt.Errorf("key %s does not exist", c.Key)
		return nil
	}
	//println("val getcommand:", val)

	c.Response <- val
	close(c.Response)

	return nil
}

func (c CommandDelete) Apply(a *Actor) error {
	_, exists := a.store[c.Key]
	if !exists {
		return fmt.Errorf("key %s does not exist", c.Key)
	}
	//println("delete command key:", c.Key)

	delete(a.store, c.Key)
	return nil
}

// Actor

type Actor struct {
	store    map[string]string
	commands chan Command
}

func NewActor() *Actor {
	a := &Actor{
		store:    make(map[string]string),
		commands: make(chan Command),
	}

	go a.run()

	return a
}

func (a *Actor) run() {
	for cmd := range a.commands {
		//fmt.Println("command:", cmd)
		err := cmd.Apply(a)
		if err != nil {
			log.Printf("Error applying command: %v", err)
		}
	}
}

func (a *Actor) Set(key, value string) error {
	errChan := make(chan error)
	a.commands <- CommandSet{Key: key, Value: value, Error: errChan}

	//fmt.Println("Set command key:", key, "value:", value)

	go func() {
		err := <-errChan
		if err != nil {
			log.Printf("Error setting key %s: %v", key, err)
		}
	}()

	return nil
}
func (a *Actor) Get(key string) (string, error) {
	responseChan := make(chan string)
	errorChan := make(chan error)
	a.commands <- CommandGet{Key: key, Response: responseChan, Error: errorChan}

	//fmt.Println("get command key", key)

	select {
	case res := <-responseChan:
		return res, nil
	case err := <-errorChan:
		return "", err
	}
}

func (a *Actor) Delete(key string) error {
	errChan := make(chan error)
	a.commands <- CommandDelete{Key: key, Error: errChan}

	go func() {
		err := <-errChan
		if err != nil {
			log.Printf("Error deleting key %s: %v", key, err)
		}
	}()

	return nil
}

func (a *Actor) List() map[string]string {
	return a.store
}

func KvHandler(a *Actor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Path[1:]
		//println("Key:", key)
		//println("request method:", r.Method)

		switch r.Method {
		case "GET":
			if key == "list" {
				keyValuePairs := a.List()
				result, err := json.Marshal(keyValuePairs)
				if err != nil {
					http.Error(w, "Error marshalling list", http.StatusInternalServerError)
					return
				}

				w.Write(result)
				return
			}

			val, err := a.Get(key)
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}

			w.Write([]byte(val))
		case "POST":
			postData, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error reading request body", http.StatusInternalServerError)
				return
			}

			fmt.Println("post data:", string(postData))
			val := string(postData)

			go a.Set(key, val)
			w.Write([]byte("Successfully set key " + key + " to value " + val))
		case "DELETE":
			go a.Delete(key)
			w.Write([]byte("Removed key " + key))
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
