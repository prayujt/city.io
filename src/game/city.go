package game

import (
	"api/database"
	"encoding/json"

	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type City struct {
	CityId     string `database:"city_id" json:"cityId"`
	Population int    `database:"population" json:"population"`
	CityName   string `database:"city_name" json:"cityName"`
	// CityOwner  string `database:"city_owner" json:"cityOwner"`
}

type Buildings struct {
	IsOwner   bool        `json:"isOwner"`
	Buildings interface{} `json:"buildings"`
}

type Building struct {
	BuildingType  string `database:"building_type"`
	BuildingLevel int    `database:"building_level"`
	BuildingName  string `database:"building_name"`
	// CityId        string `database:"city_id"`
	CityRow    int `database:"city_row"`
	CityColumn int `database:"city_column"`
	// BuildingProduction float64 `database:"building_production"`
	// HappinessChange    float64 `database:"happiness_change"`
	// BuildCost          float64 `database:"build_cost"`
	// BuildTime          int     `database:"build_time"`
}

func HandleCityRoutes(r *mux.Router) {
	r.HandleFunc("/cities/{session_id}", getCity).Methods("GET")
	r.HandleFunc("/cities/{session_id}/buildings", getBuildings).Methods("GET")

	r.HandleFunc("/cities/{session_id}/createBuilding", createBuilding).Methods("POST")
	r.HandleFunc("/cities/{session_id}/upgradeBuilding", upgradeBuilding).Methods("POST")
}

func getCity(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	sessionId := vars["session_id"]

	var city City

	defer func() {
		json.NewEncoder(response).Encode(city)
	}()

	var result []City
	database.Query(fmt.Sprintf("SELECT city_id, population, city_name FROM Cities NATURAL JOIN Sessions WHERE session_id='%s'", sessionId), &result)

	if len(result) > 0 {
		city = result[0]
	}
}

func getBuildings(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	sessionId := vars["session_id"]

	var buildings interface{}
	defer func() {
		json.NewEncoder(response).Encode(buildings)
	}()

	param := request.URL.Query()["username"]

	var query string
	if len(param) > 0 {
		query = fmt.Sprintf("SELECT building_type, building_level, building_name, city_row, city_column FROM Buildings NATURAL JOIN Cities JOIN Accounts ON city_owner=player_id WHERE username='%s';", param[0])

	} else {
		query = fmt.Sprintf("SELECT building_type, building_level, building_name, city_row, city_column FROM Buildings NATURAL JOIN Cities JOIN Sessions ON city_owner=player_id WHERE session_id='%s';", sessionId)
	}

	var buildingResult []Building
	database.Query(query, &buildingResult)

	buildings = Buildings{IsOwner: len(param) == 0, Buildings: buildingResult}
}

func createBuilding(response http.ResponseWriter, request *http.Request) {

}

func upgradeBuilding(response http.ResponseWriter, request *http.Request) {

}
