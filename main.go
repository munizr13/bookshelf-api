package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/munizr13/bookshelf-api/internal/handlers"
	"github.com/munizr13/bookshelf-api/internal/store"
)

func main() {
	s := store.New()
	h := handlers.New(s)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /books", h.ListBooks)
	mux.HandleFunc("POST /books", h.AddBook)
	mux.HandleFunc("GET /books/{id}", h.GetBook)
	mux.HandleFunc("DELETE /books/{id}", h.DeleteBook)

	addr := ":8080"
	fmt.Printf("Bookshelf API listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
