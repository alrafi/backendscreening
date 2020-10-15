package controller

import (
	"backendscreening/config/database"
	"backendscreening/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// GetDiaries is method to get all diaries
func GetDiaries(w http.ResponseWriter, r *http.Request) {
	var diaries model.Diary
	var arrDiary []model.Diary
	var response model.Response

	db := database.Connect()
	defer db.Close()

	rows, err := db.Query("Select diary_id,title,content,date from diaries")
	if err != nil {
		log.Print(err)
	}

	fmt.Println("Successfully connected to database!")

	for rows.Next() {
		if err := rows.Scan(&diaries.ID, &diaries.Title, &diaries.Content, &diaries.Date); err != nil {
			log.Fatal(err.Error())

		} else {
			arrDiary = append(arrDiary, diaries)
		}
	}

	response.Status = 200
	response.Message = "Success"
	response.Data = arrDiary

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
