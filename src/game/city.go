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
	IsOwner   bool       `json:"isOwner"`
	Buildings []Building `json:"buildings"`
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

type Status struct {
	Status bool `json:"status"`
}

func HandleCityRoutes(r *mux.Router) {
	r.HandleFunc("/cities/{session_id}", getCity).Methods("GET")
	r.HandleFunc("/cities/{session_id}/buildings", getBuildings).Methods("GET")
	r.HandleFunc("/cities/{session_id}/buildings/{city_row}/{city_column}", getBuilding).Methods("GET")

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

	var buildings Buildings
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

func getBuilding(response http.ResponseWriter, request *http.Request) {
	// vars := mux.Vars(request)
	// sessionId := vars["session_id"]
	// cityRow := vars["city_row"]
	// cityColumn := vars["city_column"]

	// param := request.URL.Query()["username"]

	// var query string
	// if len(param) > 0 {
	// 	query = fmt.Sprintf("SELECT building_type, building_level, building_name, city_row, city_column FROM Buildings NATURAL JOIN Cities JOIN Accounts ON city_owner=player_id WHERE username='%s';", param[0])

	// } else {
	// 	query = fmt.Sprintf("SELECT building_type, building_level, building_name, city_row, city_column FROM Buildings NATURAL JOIN Cities JOIN Sessions ON city_owner=player_id WHERE session_id='%s';", sessionId)
	// }
}

func createBuilding(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	sessionId := vars["session_id"]
	status := false

	defer func() {
		json.NewEncoder(response).Encode(Status{Status: status})
	}()

	var building Building
	err := json.NewDecoder(request.Body).Decode(&building)

	if err != nil {
		return
	}

	result, err := database.Execute(
		fmt.Sprintf(
			"UPDATE Accounts SET balance = balance - (SELECT build_cost FROM Building_Info WHERE building_type='%s' AND building_level=1) WHERE player_id=(SELECT player_id FROM Sessions WHERE session_id='%s')", building.BuildingType, sessionId))
	if err != nil {
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return
	}

	result, err = database.Execute(
		fmt.Sprintf(
			"INSERT INTO Buildings SELECT '%s', '%s', 1, city_id, '%d', '%d' FROM Sessions JOIN Cities ON player_id=city_owner WHERE session_id='%s'", building.BuildingName, building.BuildingType, building.CityRow, building.CityColumn, sessionId))
	if err != nil {
		return
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return
	}

	status = true
}

func upgradeBuilding(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	sessionId := vars["session_id"]
	status := false

	defer func() {
		json.NewEncoder(response).Encode(Status{Status: status})
	}()

	var building Building
	err := json.NewDecoder(request.Body).Decode(&building)

	if err != nil {
		return
	}

	result, err := database.Execute(
		fmt.Sprintf(
			"UPDATE Accounts SET balance = balance - (SELECT build_cost FROM Building_Info WHERE building_type='%s' AND building_level=(SELECT building_level + 1 FROM Buildings WHERE player_id=(SELECT player_id FROM Sessions WHERE session_id='%s') AND city_row='%d' AND city_column='%d')) WHERE player_id=(SELECT player_id FROM Sessions WHERE session_id='%s')", building.BuildingType, sessionId, building.CityRow, building.CityColumn, sessionId))
	if err != nil {
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return
	}

	result, err = database.Execute(
		fmt.Sprintf(
			"UPDATE Buildings SET building_level=building_level+1 WHERE city_id=(SELECT city_id FROM Sessions JOIN Cities ON player_id=city_owner WHERE session_id='%s') AND city_row='%d' AND city_column='%d'", sessionId, building.CityRow, building.CityColumn))
	if err != nil {
		return
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return
	}

	status = true
}
