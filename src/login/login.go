package login

import (
	"api/database"

	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	gonanoid "github.com/matoous/go-nanoid/v2"
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

func getSession(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	var expired bool
	err := database.QueryValue(fmt.Sprintf("SELECT expires_on < NOW() FROM Sessions WHERE session_id='%s'", vars["session_id"]), &expired)

	if err != nil {
		fmt.Fprintf(response, fmt.Sprintf("%v", false))
		return
	}

	fmt.Fprintf(response, fmt.Sprintf("%v", !expired))
}

func exitSession(response http.ResponseWriter, request *http.Request) {
	var session Session
	err := json.NewDecoder(request.Body).Decode(&session)

	if err != nil {
		fmt.Fprintf(response, fmt.Sprintf("%v", false))
	}

	_, err = database.Execute(fmt.Sprintf("DELETE FROM Sessions WHERE session_id='%s'", session.SessionId))

	if err != nil {
		fmt.Fprintf(response, fmt.Sprintf("%v", false))
	} else {
		fmt.Fprintf(response, fmt.Sprintf("%v", true))
	}

}

func HandleLoginRoutes(r *mux.Router) {
	r.HandleFunc("/login/createAccount", createAccount).Methods("POST")
	r.HandleFunc("/login/createSession", createSession).Methods("POST")

	r.HandleFunc("/sessions/{session_id}", getSession).Methods("GET")
	r.HandleFunc("/sessions/logout", exitSession).Methods("POST")
}
