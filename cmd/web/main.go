package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", home)
	mux.HandleFunc("GET /paste/view/{id}", pasteView)
	mux.HandleFunc("GET /paste/create", pasteCreate)
	mux.HandleFunc("POST /paste/create", pasteCreatePost)

	log.Print("starting server on :4000")

	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
