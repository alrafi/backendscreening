package controller

import (
	"backendscreening/config/database"
	"backendscreening/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Register is a method to register new user
func Register(w http.ResponseWriter, r *http.Request) {
	// var user model.User
	// var arrUser []model.User
	var response model.ResponseRegister

	db := database.Connect()
	defer db.Close()

	err := r.ParseMultipartForm(4096)
	if err != nil {
		panic(err)
	}

	username := r.FormValue("username")
	fullname := r.FormValue("fullname")
	email := r.FormValue("email")
	password := r.FormValue("password")
	birthday := r.FormValue("birthday")

	hash, errMes := bcrypt.GenerateFromPassword([]byte(password), 5)

	if errMes != nil {
		response.Message = "Error While Hashing Password, Try Again"
		json.NewEncoder(w).Encode(response)
		return
	}
	password = string(hash)

	_, err = db.Exec("INSERT INTO users (username, fullname, email, password, birthday) values ($1,$2,$3,$4,$5)", username, fullname, email, password, birthday)

	if err != nil {
		log.Print(err)
		log.Print("after db exec")
		return
	}

	response.Status = 200
	response.Message = "Success Register New User"
	// response.Data = model.User{Username: username, Fullname: fullname, Email: email, Password: password, Birthday: birthday}
	log.Print("Insert data to database")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	return
}

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
