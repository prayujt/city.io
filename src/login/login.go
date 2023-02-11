package login

import (
	"api/database"

	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/matoous/go-nanoid/v2"
)

type Account struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Session struct {
	Status    bool   `json:"status"`
	SessionId string `json:"sessionId"`
}

func createAccount(response http.ResponseWriter, request *http.Request) {
	log.Println("Received request to /login/createAccount")

	status := false
	defer func() {
		fmt.Fprintf(response, fmt.Sprintf("%v", status))
	}()

	var acc Account
	err := json.NewDecoder(request.Body).Decode(&acc)

	if err != nil {
		return
	}

	_, err = database.Execute(fmt.Sprintf("INSERT INTO Accounts (player_id, username, password) VALUES (uuid(), '%s', '%s')", acc.Username, acc.Password))
	if err == nil {
		status = true
	}
}

func createSession(response http.ResponseWriter, request *http.Request) {
	log.Println("Received request to /login/createSession")

	status := false
	nullSession := Session{Status: false, SessionId: ""}

	sessionId, err := gonanoid.New()
	session := Session{Status: true, SessionId: sessionId}

	defer func() {
		if status {
			json.NewEncoder(response).Encode(session)
		} else {
			json.NewEncoder(response).Encode(nullSession)
		}
	}()

	if err != nil {
		return
	}

	var acc Account
	err = json.NewDecoder(request.Body).Decode(&acc)

	if err != nil {
		return
	}

	result, err := database.Execute(
		fmt.Sprintf(
			"INSERT INTO Sessions (session_id, player_id, expires_on) SELECT '%s', player_id, DATE_ADD(NOW(), INTERVAL 24 HOUR) FROM Accounts WHERE username='%s' AND password='%s';", sessionId, acc.Username, acc.Password))

	if err != nil {
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return
	}

	status = true
}

func HandleLoginRoutes(r *mux.Router) {
	r.HandleFunc("/login/createAccount", createAccount).Methods("POST")
	r.HandleFunc("/login/createSession", createSession).Methods("POST")
}
