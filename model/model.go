package model

// User is model of user
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	Birthday string `json:"birthday"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Diary is model of diary
type Diary struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Date    string `json:"time"`
}

// Response is model of response
type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []Diary
}

// ResponseRegister is model of response result
type ResponseRegister struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    User
	Token   string `json:"token"`
}

// ResponseLogin is model of response result
type ResponseLogin struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    User
	Token   string `json:"token"`
}

// ResponseUsers is model of response
type ResponseUsers struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []User
}
