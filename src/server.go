package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}

	dbUrl := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_DB_NAME"))

	_, err = sql.Open("mysql", dbUrl)

	if err != nil {
		panic(err)
	}

	fmt.Println("Serving at localhost:8000...")
	http.ListenAndServe(":8000", nil)
}
