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
	Username string
	Password string
}

type AccountMatch struct {
	Status bool
	Uuid   string
}

func createAccount(response http.ResponseWriter, request *http.Request) {
	log.Println("Received request to /login/createAccount")
	var acc Account
	err := json.NewDecoder(request.Body).Decode(&acc)

	if err != nil {
		panic(err)
	}

	err = execute(fmt.Sprintf("INSERT INTO Accounts (uuid, username, password) VALUES (uuid(), '%s', '%s')", acc.Username, acc.Password))
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

	var uuid string
	err = queryValue(fmt.Sprintf("SELECT uuid FROM Accounts WHERE username='%s' AND password='%s'", acc.Username, acc.Password), &uuid)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	status := AccountMatch{Status: true, Uuid: uuid}
	if uuid == "" {
		status.Status = false
	}

	json.NewEncoder(response).Encode(status)
}

func handleLoginRoutes(r *mux.Router) {
	r.HandleFunc("/login/createAccount", createAccount).Methods("POST")
	r.HandleFunc("/login/verifyAccount", verifyAccount).Methods("POST")
}
