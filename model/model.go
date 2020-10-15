package model

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
