package main

import (
	"backendscreening/controller"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func home(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("welcome")
}

// init diaries
// var diaries []Diary

func main() {
	r := mux.NewRouter()

	// get current time
	// currentTime := time.Now()

	// mock data
	// diaries = append(diaries, Diary{ID: "1", Title: "Day One as Software Engineer Facebook", Content: "I hope I can give full of my skills to this company", Date: currentTime.Format("2006-01-02")})

	r.HandleFunc("/", home).Methods("GET")
	r.HandleFunc("/api/diaries", controller.GetDiaries).Methods("GET")

	log.Fatal(http.ListenAndServe(":8000", r))
	fmt.Println("Successfully connected to port 8000")
}
