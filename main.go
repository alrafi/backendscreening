package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func home(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("welcome")
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", home).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", r))
}
