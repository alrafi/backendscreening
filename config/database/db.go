package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func getEnv() {
	env := godotenv.Load("credentials.env")
	if env != nil {
		log.Fatal("Error on loading .env file")
	}
}

// Connect is method to connect with postgres db
func Connect() *sql.DB {
	getEnv()
	var (
		host     = os.Getenv("DATABASE_HOST")
		port     = os.Getenv("DATABASE_PORT")
		user     = os.Getenv("DATABASE_USER")
		password = os.Getenv("DATABASE_PASSWORD")
		dbname   = os.Getenv("DATABASE_NAME")
	)
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal(err)
	}

	return db
}
