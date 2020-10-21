package controller

import (
	"backendscreening/config/database"
	"backendscreening/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// Register is a method to register new user
func Register(w http.ResponseWriter, r *http.Request) {
	var user model.User
	// var arrUser []model.User
	var response model.ResponseRegister

	db := database.Connect()
	defer db.Close()

	// err := r.ParseMultipartForm(4096)
	// if err != nil {
	// 	panic(err)
	// }

	// username := r.FormValue("username")
	// fullname := r.FormValue("fullname")
	// email := r.FormValue("email")
	// password := r.FormValue("password")
	// birthday := r.FormValue("birthday")

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	sign := jwt.New(jwt.GetSigningMethod("HS256"))
	claims := sign.Claims.(jwt.MapClaims)
	claims["user"] = user.Username
	token, err := sign.SignedString([]byte("secret"))

	password := &user.Password

	hash, errMes := bcrypt.GenerateFromPassword([]byte(*password), 5)

	if errMes != nil {
		response.Message = "Error While Hashing Password, Try Again"
		json.NewEncoder(w).Encode(response)
		return
	}
	*password = string(hash)

	_, err = db.Exec("INSERT INTO users (username, fullname, email, password, birthday) values ($1,$2,$3,$4,$5)", user.Username, user.Fullname, user.Email, *password, user.Birthday)

	if err != nil {
		log.Print(err)
		log.Print("after db exec")
		return
	}

	// arrUser = append(arrUser, user)

	response.Status = 200
	response.Message = "Success Register New User"
	response.Data = user
	response.Token = token
	// response.Data = model.User{Username: username, Fullname: fullname, Email: email, Password: password, Birthday: birthday}
	log.Print("Insert data to database")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	return
}

// Login is method to log in
func Login(w http.ResponseWriter, r *http.Request) {
	var user model.User
	// var arrUser []model.User
	var response model.ResponseLogin

	db := database.Connect()
	defer db.Close()

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	plainPass := []byte(user.Password)
	// plainPass := []byte("halo")
	// log.Print(plainPass)

	// log.Print(user)

	row := db.QueryRow("SELECT username, fullname, password, email, birthday FROM users WHERE username=$1", user.Username)
	err = row.Scan(&user.Username, &user.Fullname, &user.Password, &user.Email, &user.Birthday)

	if err != nil {
		log.Print(err)
		log.Print("ada error")
	}

	// log.Print(plainPass)
	// log.Print(user.Password)
	hashPass := []byte(user.Password)
	err = bcrypt.CompareHashAndPassword(hashPass, plainPass)
	if err != nil {
		log.Println(err)
		log.Print("err pass")
	}

	sign := jwt.New(jwt.GetSigningMethod("HS256"))
	claims := sign.Claims.(jwt.MapClaims)
	claims["user"] = user.Username
	token, err := sign.SignedString([]byte("secret"))

	response.Status = 200
	response.Message = "Success Login"
	response.Data = user
	response.Token = token

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetUsers is method to get all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	var user model.User
	var arrUser []model.User
	var response model.ResponseUsers

	db := database.Connect()
	defer db.Close()

	rows, err := db.Query("Select username, fullname, email, birthday from users")
	if err != nil {
		log.Print(err)
	}

	for rows.Next() {
		if err := rows.Scan(&user.Username, &user.Fullname, &user.Email, &user.Birthday); err != nil {
			log.Fatal(err.Error())

		} else {
			arrUser = append(arrUser, user)
		}
	}

	response.Status = 200
	response.Message = "Success"
	response.Data = arrUser

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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

// JwtVerify is auth midlleware
func JwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var header = r.Header.Get("Authorization") //Grab the token from the header
		// var response model.ResponseLogin

		// header = strings.TrimSpace(header)

		if header == "" {
			//Token is missing, returns with error code 403 Unauthorized
			w.WriteHeader(http.StatusForbidden)
			// response.Message = "Missing auth token"
			// json.NewEncoder(w).Encode(response)
			return
		}

		token, err := jwt.Parse(header, func(token *jwt.Token) (interface{}, error) {
			if jwt.GetSigningMethod("HS256") != token.Method {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte("secret"), nil
		})

		if token != nil && err == nil {
			fmt.Println("token verified")
			next.ServeHTTP(w, r)
		} else {
			fmt.Println("error")
			json.NewEncoder(w).Encode("Error token authentication")
		}
	})
}
