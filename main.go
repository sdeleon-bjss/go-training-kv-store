package main

import (
	"github.com/sdeleon-bjss/store"
	"net/http"
	"os"
	"os/signal"
)

func main() {
	kva := store.NewActor()
	println("kv store created")

	go kva.Set("example_key", "example value")

	http.HandleFunc("/", store.KvHandler(kva))

	server := &http.Server{
		Addr: ":8080",
	}

	go func() {
		println("starting server - listeining on port 8080")

		err := server.ListenAndServe()
		if err != nil {
			return
		}
	}()

	// shutdown and cleanup

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	kva.Close()

	err := server.Shutdown(nil)
	if err != nil {
		return
	}

	println("server stopped")
}

// Notes learned from shutdown and cleanup:
// created a channel of type os.Signal buffered space for 1
// - an os.Interrupt signal is like (ctrl+c)
//
// wait for signal and block main thread
// close the actor safely
// 	- this is a way to signal the actor to stop processing commands because we are shutting down th eserver
