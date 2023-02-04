package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
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
	fmt.Println("Received request to /login/createAccount")
	var acc Account
	err := json.NewDecoder(request.Body).Decode(&acc)

	if err != nil {
		panic(err)
	}

	_, err = db.Exec("INSERT INTO Accounts (uuid, username, password) VALUES (uuid(), ?, ?)", acc.Username, acc.Password)
	if err != nil {
		fmt.Fprintf(response, "false")
	} else {
		fmt.Fprintf(response, "true")
	}
}

func verifyAccount(response http.ResponseWriter, request *http.Request) {
	fmt.Println("Received request to /login/verifyAccount")
	var acc Account
	err := json.NewDecoder(request.Body).Decode(&acc)

	if err != nil {
		panic(err)
	}

	var uuid string
	err = db.QueryRow("SELECT uuid FROM Accounts WHERE username=? and password=?", acc.Username, acc.Password).Scan(&uuid)

	status := AccountMatch{Status: true, Uuid: uuid}

	switch err {
	case sql.ErrNoRows:
		status.Status = false
		json.NewEncoder(response).Encode(status)
	case nil:
		json.NewEncoder(response).Encode(status)
	default:
		panic(err)
	}
}

func handleLoginRoutes(r *mux.Router) {
	r.HandleFunc("/login/createAccount", createAccount).Methods("POST")
	r.HandleFunc("/login/verifyAccount", verifyAccount).Methods("POST")
}
