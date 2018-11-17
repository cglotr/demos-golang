package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
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
var signingKey = []byte("secret")

func main() {
	r := mux.NewRouter()

	r.Handle("/", http.FileServer(http.Dir("./views/")))
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	r.Handle("/get-token", getTokenHandler).Methods("GET")
	r.Handle("/products", jwtHandler.Handler(productsHandler)).Methods("GET")
	r.Handle(
		"/products/{slug}/feedback",
		jwtHandler.Handler(postProductFeedbackHandler),
	).Methods("POST")
	r.Handle("/status", statusHandler).Methods("GET")

	http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, r))
}

var getTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "Faithes"

	tokenString, _ := token.SignedString(signingKey)

	w.Write([]byte(tokenString))
})

var jwtHandler = jwtmiddleware.New(jwtmiddleware.Options{
	SigningMethod: jwt.SigningMethodHS256,
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	},
})

var notImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
})

var postProductFeedbackHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var product product
	vars := mux.Vars(r)
	slug := vars["slug"]

	for _, p := range products {
		if p.Slug == slug {
			product = p
		}
	}

	w.Header().Set("Content-Type", "application/json")

	if product.Slug != "" {
		payload, _ := json.Marshal(product)
		w.Write([]byte(payload))
	} else {
		w.Write([]byte("Product Not Found"))
	}
})

var productsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	payload, _ := json.Marshal(products)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
})

var statusHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("API server is up & running"))
})
