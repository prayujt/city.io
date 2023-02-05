package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

var db *sql.DB

func query[T any](query_ string, arr *[]T) {
	rows, err := db.Query(query_)
	if err != nil {
		panic(err)
	}

	structType := reflect.TypeOf(arr).Elem().Elem()

	var fields []string
	cols, err := rows.Columns()
	for _, col := range cols {
		found := false
		for i := 0; i < structType.NumField(); i++ {
			if strings.EqualFold(structType.Field(i).Name, col) || strings.EqualFold(structType.Field(i).Tag.Get("json"), col) {
				fields = append(fields, structType.Field(i).Name)
				found = true
			}
		}
		if !found {
			fields = append(fields, "")
		}
	}

	for rows.Next() {
		newStruct := reflect.New(structType).Elem()
		var properties []interface{}

		for _, field := range fields {
			if field == "" {
				var nothing interface{}
				properties = append(properties, &nothing)
			} else {
				properties = append(properties, newStruct.FieldByName(field).Addr().Interface())
			}
		}

		rows.Scan(properties...)
		*arr = append(*arr, newStruct.Interface().(T))
	}
}

func queryValue[T any](query_ string, value *T) error {
	err := db.QueryRow(query_).Scan(value)
	return err
}

func execute(exec string) error {
	_, err := db.Exec(exec)
	return err
}

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbUrl := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_DB_NAME"))

	db, err = sql.Open("mysql", dbUrl)

	if err != nil {
		log.Fatal("Couldn't connect to db")
	}

	log.Println(fmt.Sprintf("Serving at %s:%s...", os.Getenv("API_HOST"), os.Getenv("API_PORT")))
	router := mux.NewRouter()

	// include other file routes here, passing in the router
	handleLoginRoutes(router)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	server := &http.Server{
		Handler:      handler,
		Addr:         fmt.Sprintf("0.0.0.0:%s", os.Getenv("API_PORT")),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
