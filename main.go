package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Diary struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Date    string `json:"time"`
}

func home(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("welcome")
}

func getDiaries(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	json.NewEncoder(w).Encode(diaries)
}

// init diaries
var diaries []Diary

func main() {
	r := mux.NewRouter()

	// get current time
	currentTime := time.Now()

	// mock data
	diaries = append(diaries, Diary{ID: "1", Title: "Day One as Software Engineer Facebook", Content: "I hope I can give full of my skills to this company", Date: currentTime.Format("2006-01-02")})

	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/api/diaries/{id}", getDiaries).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", r))
}
