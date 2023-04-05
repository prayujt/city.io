package tests

import (
	"api/database"
	"api/game"
	"api/login"

	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var URL string
var jwtToken string

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
	game.HandleCityRoutes(router)
	game.HandleVisitRoutes(router)
	game.HandleArmyRoutes(router)

	server := httptest.NewServer(router)
	URL = server.URL

	acc := login.Account{
		Username: "root",
		Password: "root",
	}

	Post("/login/createAccount", acc)

	response := Post("/login/createSession", acc)

	var session login.JWT
	json.Unmarshal(response, &session)

	if session.Token == "" {
		fmt.Println("Failed to initialize token for tests")
		return
	}

	city := game.CityNameChange{
		CityNameNew: "monkee monkee",
	}

	response1 := Post("/cities/updateName", city, session.Token)
	var result game.Status

	json.Unmarshal(response1, &result)

	if !result.Status {
		fmt.Println("youre gay")
	}

	jwtToken = session.Token

	acc2 := login.Account{
		Username: "player1",
		Password: "player1",
	}

	Post("/login/createAccount", acc2)

	response2 := Post("/login/createSession", acc2)

	var session1 login.JWT
	json.Unmarshal(response2, &session1)

	if session1.Token == "" {
		fmt.Println("Failed to initialize token for tests")
		return
	}

	city1 := game.CityNameChange{
		CityNameNew: "monkee",
	}

	response3 := Post("/cities/updateName", city1, session1.Token)
	var result1 game.Status

	json.Unmarshal(response3, &result1)

	m.Run()
}

func Post(endpoint string, reqBody any, token ...string) []byte {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(reqBody)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", URL+endpoint, &buf)

	if err != nil {
		panic(err)
	}

	if len(token) > 0 {
		req.Header.Set("Token", token[0])
	} else if jwtToken != "" {
		req.Header.Set("Token", jwtToken)
	}

	res, err := client.Do(req)
	resBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	return resBody
}

func Get(endpoint string, token ...string) []byte {
	client := &http.Client{}

	req, err := http.NewRequest("GET", URL+endpoint, nil)

	if err != nil {
		panic(err)
	}

	if len(token) > 0 {
		req.Header.Set("Token", token[0])
	} else if jwtToken != "" {
		req.Header.Set("Token", jwtToken)
	}

	res, err := client.Do(req)
	resBody, err := ioutil.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	return resBody
}
