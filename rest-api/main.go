package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

var nextID = 1
var spendings []spending

type spending struct {
	ID   string `json:"id,omitempty"`
	Item string `json:"item,omitempty"`
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/spendings", getSpendings).Methods("GET")
	router.HandleFunc("/spendings", createSpending).Methods("POST")
	router.HandleFunc("/spendings/{id}", getSpending).Methods("GET")
	router.HandleFunc("/spendings/{id}", deleteSpending).Methods("DELETE")
	router.HandleFunc("/spendings/{id}", updateSpending).Methods("PATCH")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func createSpending(w http.ResponseWriter, r *http.Request) {
	var spending spending
	_ = json.NewDecoder(r.Body).Decode(&spending)
	spending.ID = strconv.Itoa(nextID)
	nextID++
	spendings = append(spendings, spending)
	json.NewEncoder(w).Encode(spending)
}

func deleteSpending(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var index = -1

	for i, item := range spendings {
		if item.ID == params["id"] {
			index = i
		}
	}

	if index != -1 {
		spendings = append(spendings[:index], spendings[index+1:]...)
	}

	json.NewEncoder(w).Encode(spendings)
}

func getSpending(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range spendings {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&spending{})
}

func getSpendings(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(spendings)
}

func updateSpending(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var spending spending
	json.NewDecoder(r.Body).Decode(&spending)

	id := params["id"]
	var itemUpdate = spending.Item

	for index, item := range spendings {
		if item.ID == id {
			spendings[index].Item = itemUpdate
			json.NewEncoder(w).Encode(item)
			return
		}
	}

	json.NewEncoder(w)
}
