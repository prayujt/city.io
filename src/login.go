package main

import (
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

func handleLoginRoutes(r *mux.Router) {
	r.HandleFunc("/login/createAccount", createAccount).Methods("POST")
}
