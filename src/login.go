package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Account struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AccountMatch struct {
	Status   bool   `json:"status"`
	PlayerId string `json:"playerId"`
}

func createAccount(response http.ResponseWriter, request *http.Request) {
	log.Println("Received request to /login/createAccount")
	var acc Account
	err := json.NewDecoder(request.Body).Decode(&acc)

	if err != nil {
		panic(err)
	}

	err = execute(fmt.Sprintf("INSERT INTO Accounts (player_id, username, password) VALUES (uuid(), '%s', '%s')", acc.Username, acc.Password))
	if err != nil {
		fmt.Fprintf(response, "false")
	} else {
		fmt.Fprintf(response, "true")
	}
}

func verifyAccount(response http.ResponseWriter, request *http.Request) {
	log.Println("Received request to /login/verifyAccount")
	var acc Account
	err := json.NewDecoder(request.Body).Decode(&acc)

	if err != nil {
		panic(err)
	}

	var playerId string
	err = queryValue(fmt.Sprintf("SELECT player_id FROM Accounts WHERE username='%s' AND password='%s'", acc.Username, acc.Password), &playerId)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	status := AccountMatch{Status: playerId != "", PlayerId: playerId}

	json.NewEncoder(response).Encode(status)
}

func handleLoginRoutes(r *mux.Router) {
	r.HandleFunc("/login/createAccount", createAccount).Methods("POST")
	r.HandleFunc("/login/verifyAccount", verifyAccount).Methods("POST")
}
