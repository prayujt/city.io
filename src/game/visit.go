package game

import (
	"api/database"
	"encoding/json"

	"net/http"

	"github.com/gorilla/mux"
)

type Ownership struct {
	CityName  string `database:"city_name" json:"cityName"`
	CityOwner string `database:"username" json:"cityOwner"`
}

type Player struct {
	Username string  `database:"username" json:"username"`
	Balance  float64 `database:"balance" json:"balance"`
}

func HandleVisitRoutes(r *mux.Router) {
	r.HandleFunc("/cities", getCityList).Methods("GET")
	r.HandleFunc("/leaderboard", getLeaderBoard).Methods("GET")
	r.HandleFunc("/towns", getTownList).Methods("GET")
}

func getCityList(response http.ResponseWriter, request *http.Request) {
	var list []Ownership
	database.Query("SELECT city_name, username FROM Cities JOIN Accounts ON city_owner=player_id WHERE town=0;", &list)
	json.NewEncoder(response).Encode(list)
}

func getTownList(response http.ResponseWriter, request *http.Request) {
	var list []Ownership
	database.Query("SELECT city_name, username FROM Cities JOIN Accounts ON city_owner=player_id WHERE town=1;", &list)
	json.NewEncoder(response).Encode(list)
}

func getLeaderBoard(response http.ResponseWriter, request *http.Request) {
	var list []Player
	database.Query("SELECT username, balance FROM Accounts ORDER BY balance DESC;", &list)
	json.NewEncoder(response).Encode(list)
}
