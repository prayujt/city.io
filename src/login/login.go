package login

import (
	"api/auth"
	"api/database"

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

type JWT struct {
	Token string `json:"token"`
}

type Status struct {
	Status bool `json:"status"`
}

func HandleLoginRoutes(r *mux.Router) {
	r.HandleFunc("/login/createAccount", createAccount).Methods("POST")
	r.HandleFunc("/login/createSession", createSession).Methods("POST")

	r.HandleFunc("/sessions/validate", validateSession).Methods("GET")
}

func createAccount(response http.ResponseWriter, request *http.Request) {
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
	result, err := database.Execute(
		fmt.Sprintf(
			"INSERT INTO Buildings (building_type, building_level, city_id, city_row, city_column) SELECT 'City Hall', 1, city_id, 4, 6 FROM Cities WHERE city_owner=(SELECT player_id FROM Accounts where username='%s');", acc.Username))

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
	var token string

	defer func() {
		json.NewEncoder(response).Encode(JWT{Token: token})
	}()

	var acc Account
	err := json.NewDecoder(request.Body).Decode(&acc)

	if err != nil {
		log.Println(err)
		response.WriteHeader(400)
		return
	}

	type Player struct {
		PlayerId   string `database:"player_id"`
		Authorized bool   `database:"authorized"`
	}

	var player []Player

	database.Query(
		fmt.Sprintf("SELECT player_id, password=SHA2('%s', 256) AS authorized FROM Accounts WHERE username='%s'", acc.Password, acc.Username), &player)

	if len(player) == 0 {
		response.WriteHeader(401)
		return
	}

	if !player[0].Authorized {
		response.WriteHeader(401)
		return
	}

	token, err = auth.GenerateJWT(acc.Username, player[0].PlayerId)

	if err != nil {
		log.Println(err)
		response.WriteHeader(400)
		return
	}
}

func validateSession(response http.ResponseWriter, request *http.Request) {
	var valid bool = false

	defer func() {
		json.NewEncoder(response).Encode(Status{Status: valid})
	}()

	if request.Header["Token"] != nil {
		_, err := auth.ParseJWT(request.Header["Token"][0])

		if err != nil {
			response.WriteHeader(400)
			return
		}
		valid = true
	}
}
