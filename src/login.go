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

func createAccount(response http.ResponseWriter, request *http.Request) {
	fmt.Println("Received request to /login/createAccount")
	var acc Account
	err := json.NewDecoder(request.Body).Decode(&acc)

	if err != nil {
		panic(err)
	}

	_, err = db.Exec("INSERT INTO Accounts (username, password) VALUES (?, ?, ?)", acc.Username, acc.Password)
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

	var username string
	err = db.QueryRow("SELECT username FROM Accounts WHERE username=? and password=?", acc.Username, acc.Password).Scan(&username)

	switch err {
	case sql.ErrNoRows:
		fmt.Fprintf(response, "false")
	case nil:
		fmt.Fprintf(response, "true")
	default:
		panic(err)
	}
}

func handleLoginRoutes(r *mux.Router) {
	r.HandleFunc("/login/createAccount", createAccount).Methods("POST")
	r.HandleFunc("/login/verifyAccount", verifyAccount).Methods("POST")
}
