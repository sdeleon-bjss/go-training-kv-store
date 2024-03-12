package main

import (
	"bjss-kv-store/store"
	"net/http"
)

func main() {
	kva := store.NewActor()
	println("kv store created")

	http.HandleFunc("/", store.KvHandler(kva))

	println("server started")
	http.ListenAndServe(":8080", nil)
	println("server stopped")

	// Todo - shutdown signal handling
	// Todo - closing kv store?/chnanels so no more commands can be sent
}
