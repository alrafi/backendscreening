package controller

import (
	"backendscreening/config/database"
	"backendscreening/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"unicode"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// Register is a method to register new user
func Register(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user model.User
	// var arrUser []model.User
	var response model.ResponseRegister

	db := database.Connect()
	defer db.Close()

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	password := user.Password
	errPass := VerifyPassword(password)

	if errPass != "" {
		log.Print(errPass)
		// log.Print("err password")
		response.Status = 400
		response.Message = errPass

		json.NewEncoder(w).Encode(response)
		return
	}

	sign := jwt.New(jwt.GetSigningMethod("HS256"))
	claims := sign.Claims.(jwt.MapClaims)
	claims["user"] = user.ID
	token, err := sign.SignedString([]byte("secret"))

	hash, errMes := bcrypt.GenerateFromPassword([]byte(password), 5)

	if errMes != nil {
		response.Message = "Error While Hashing Password, Try Again"
		json.NewEncoder(w).Encode(response)
		return
	}
	password = string(hash)

	_, err = db.Exec("INSERT INTO users (user_id, username, fullname, email, password, birthday) values (DEFAULT, $1,$2,$3,$4,$5)", user.Username, user.Fullname, user.Email, password, user.Birthday)

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

	// w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	return
}

// Login is method to log in
func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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

	var data string

	if user.Username != "" {
		log.Print("username != nil")
		data = user.Username
		row := db.QueryRow("SELECT user_id, username, fullname, password, email, birthday FROM users WHERE username=$1", data)
		err = row.Scan(&user.ID, &user.Username, &user.Fullname, &user.Password, &user.Email, &user.Birthday)
	} else if user.Email != "" {
		log.Print("email != nil")
		data = user.Email
		row := db.QueryRow("SELECT user_id, username, fullname, password, email, birthday FROM users WHERE email=$1", data)
		err = row.Scan(&user.ID, &user.Username, &user.Fullname, &user.Password, &user.Email, &user.Birthday)
	}

	if err != nil {
		log.Print(err)
		log.Print("ada error")
		response.Status = 400
		response.Message = "Error: akun tersebut tidak terdaftar"
		json.NewEncoder(w).Encode(response)
		return
	}

	hashPass := []byte(user.Password)
	err = bcrypt.CompareHashAndPassword(hashPass, plainPass)
	if err != nil {
		log.Println(err)
		log.Print("err pass")
		response.Status = 400
		response.Message = "Error: username/email dan password tidak cocok"
		json.NewEncoder(w).Encode(response)
		return
	}

	sign := jwt.New(jwt.GetSigningMethod("HS256"))
	claims := sign.Claims.(jwt.MapClaims)
	claims["user"] = user.ID
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

	rows, err := db.Query("Select user_id, username, fullname, email, birthday from users")
	if err != nil {
		log.Print(err)
	}

	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Username, &user.Fullname, &user.Email, &user.Birthday); err != nil {
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

// GetDiaries is method to get all diaries from a user
func GetDiaries(w http.ResponseWriter, r *http.Request) {
	var diaries model.Diary
	var arrDiary []model.Diary
	var response model.Response

	db := database.Connect()
	defer db.Close()

	var header = r.Header.Get("Authorization")
	// log.Print(header)
	encode := extractClaims(header)
	claims := encode.Claims.(jwt.MapClaims)
	log.Print(claims)
	selectedUser := claims["user"]
	log.Print(selectedUser)

	rows, err := db.Query("Select diary_id,title,content,date from diaries WHERE user_id=$1", selectedUser)
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

// CreateDiaries is method to create a diary
func CreateDiaries(w http.ResponseWriter, r *http.Request) {
	var diary model.Diary
	var arrDiary []model.Diary
	var response model.Response

	db := database.Connect()
	defer db.Close()

	var header = r.Header.Get("Authorization")
	err := json.NewDecoder(r.Body).Decode(&diary)
	// log.Print(header)
	encode := extractClaims(header)
	claims := encode.Claims.(jwt.MapClaims)
	log.Print(claims)
	selectedUser := claims["user"]
	log.Print(selectedUser)

	_, err = db.Exec("INSERT INTO diaries (diary_id, user_id, title, content, date) VALUES (DEFAULT, $1, $2, $3, $4)", selectedUser, diary.Title, diary.Content, diary.Date)

	if err != nil {
		log.Print(err)
		return
	}

	fmt.Println("Successfully connected to database!")

	arrDiary = append(arrDiary, diary)

	response.Status = 200
	response.Message = "Success add diary"
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
			// fmt.Println(token)
			next.ServeHTTP(w, r)
		} else {
			fmt.Println("error")
			json.NewEncoder(w).Encode("Error token authentication")
		}
	})
}

// VerifyPassword is
func VerifyPassword(password string) string {
	var uppercasePresent bool
	var lowercasePresent bool
	var numberPresent bool
	var specialCharPresent bool
	const minPassLength = 6
	const maxPassLength = 32
	var passLen int
	var errorString string

	for _, ch := range password {
		switch {
		case unicode.IsNumber(ch):
			numberPresent = true
			passLen++
		case unicode.IsUpper(ch):
			uppercasePresent = true
			passLen++
		case unicode.IsLower(ch):
			lowercasePresent = true
			passLen++
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			specialCharPresent = true
			passLen++
		case ch == ' ':
			passLen++
		}
	}
	appendError := func(err string) {
		if len(strings.TrimSpace(errorString)) != 0 {
			errorString += ", " + err
		} else {
			errorString = err
		}
	}
	if !lowercasePresent {
		appendError("lowercase letter missing")
	}
	if !uppercasePresent {
		appendError("uppercase letter missing")
	}
	if !numberPresent {
		appendError("atleast one numeric character required")
	}
	if !specialCharPresent {
		appendError("special character missing")
	}
	if !(minPassLength <= passLen && passLen <= maxPassLength) {
		appendError(fmt.Sprintf("password length must be between %d to %d characters long", minPassLength, maxPassLength))
	}

	if len(errorString) != 0 {
		return errorString
	}
	return ""
}

func extractClaims(tokenStr string) *jwt.Token {

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("secret"), nil
	})

	if token != nil && err == nil {
		return token
	} else {
		return nil
	}
}
