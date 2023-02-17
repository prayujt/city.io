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
	SessionId string `json:"sessionId"`
}

type Status struct {
	Status bool `json:"status"`
}

func HandleLoginRoutes(r *mux.Router) {
	r.HandleFunc("/login/createAccount", createAccount).Methods("POST")
	r.HandleFunc("/login/createSession", createSession).Methods("POST")

	r.HandleFunc("/sessions/{session_id}", getSession).Methods("GET")
	r.HandleFunc("/sessions/logout", exitSession).Methods("POST")
}

func createAccount(response http.ResponseWriter, request *http.Request) {
	log.Println("Received request to /login/createAccount")

	status := false
	defer func() {
		json.NewEncoder(response).Encode(Status{Status: status})
	}()

	var acc Account
	err := json.NewDecoder(request.Body).Decode(&acc)

	if err != nil {
		return
	}

	_, err = database.Execute(fmt.Sprintf("INSERT INTO Accounts (player_id, username, password) VALUES (uuid(), '%s', SHA2('%s', 256))", acc.Username, acc.Password))
	if err != nil {
		return
	}

	var city string = "City Hall"
	result, err := database.Execute(
		fmt.Sprintf(
			"INSERT INTO Buildings (building_name, building_type, building_level, city_id, city_row, city_column) SELECT '%s', '%s', 1, city_id, 4, 4 FROM Cities WHERE city_owner=(SELECT player_id FROM Accounts where username='%s');", city, city, acc.Username))

	if err != nil {
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return
	}

	status = true
}

func createSession(response http.ResponseWriter, request *http.Request) {
	log.Println("Received request to /login/createSession")

	status := false
	sessionId, err := gonanoid.New()

	defer func() {
		if !status {
			sessionId = ""
		}
		json.NewEncoder(response).Encode(Session{SessionId: sessionId})
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
			"INSERT INTO Sessions (session_id, player_id, expires_on) SELECT '%s', player_id, DATE_ADD(NOW(), INTERVAL 24 HOUR) FROM Accounts WHERE username='%s' AND password=SHA2('%s', 256);", sessionId, acc.Username, acc.Password))

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
	// log.Println("Received request to /sessions/{session_id}")
	vars := mux.Vars(request)
	expired := true

	defer func() {
		json.NewEncoder(response).Encode(Status{Status: !expired})
	}()

	_ = database.QueryValue(fmt.Sprintf("SELECT expires_on < NOW() FROM Sessions WHERE session_id='%s'", vars["session_id"]), &expired)
}

func exitSession(response http.ResponseWriter, request *http.Request) {
	log.Println("Received request to /sessions/logout")
	var session Session
	status := false

	defer func() {
		json.NewEncoder(response).Encode(Status{Status: status})
	}()

	err := json.NewDecoder(request.Body).Decode(&session)

	if err != nil {
		return
	}

	_, err = database.Execute(fmt.Sprintf("DELETE FROM Sessions WHERE session_id='%s'", session.SessionId))

	if err == nil {
		status = true
	}
}
