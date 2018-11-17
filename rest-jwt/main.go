package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type product struct {
	Description string
	ID          int
	Name        string
	Slug        string
}

var products = []product{
	product{
		Description: "Test product.",
		ID:          1,
		Name:        "Test 001",
		Slug:        "test-001",
	},
	product{
		Description: "Test product.",
		ID:          2,
		Name:        "Test 002",
		Slug:        "test-002",
	},
}

func main() {
	r := mux.NewRouter()

	r.Handle("/", http.FileServer(http.Dir("./views/")))
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	r.Handle("/products", productsHandler).Methods("GET")
	r.Handle("/products/{slug}/feedback", notImplemented).Methods("POST")
	r.Handle("/status", notImplemented).Methods("GET")

	http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, r))
}

var notImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
})

var productsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	payload, _ := json.Marshal(products)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
})
