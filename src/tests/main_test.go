package tests

import (
	"api/database"
	"api/login"

	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var URL string

func TestMain(m *testing.M) {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic(err)
	}

	dbUrl := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_TEST_DB_NAME"))
	database.InitDatabase(dbUrl)
	database.Reset()

	router := mux.NewRouter()
	login.HandleLoginRoutes(router)

	server := httptest.NewServer(router)
	URL = server.URL
	m.Run()
}

func Post(endpoint string, reqBody any) []byte {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(reqBody)
	if err != nil {
		panic(err)
	}

	res, err := http.Post(URL+endpoint, "application/json", &buf)
	if err != nil {
		panic(err)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return resBody
}

func Get(endpoint string) []byte {
	res, err := http.Get(URL + endpoint)
	if err != nil {
		panic(err)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}
	return resBody
}
