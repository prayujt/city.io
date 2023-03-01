package game

import (
	"api/database"
	"encoding/json"
	"math"

	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type City struct {
	// CityId     string `database:"city_id" json:"cityId"`
	CityName           string  `database:"city_name" json:"cityName"`
	Population         int     `database:"population" json:"population"`
	PopulationCapacity int     `database:"population_capacity" json:"populationCapacity"`
	PlayerBalance      float64 `database:"balance" json:"playerBalance"`
	CityOwner          string  `database:"username" json:"cityOwner"`
}

type CityNameChange struct {
	CityNameOriginal string `json:"cityNameOriginal"`
	CityNameNew      string `json:"cityNameNew"`
}

type Buildings struct {
	IsOwner   bool       `json:"isOwner"`
	Buildings []Building `json:"buildings"`
}

type Building struct {
	BuildingType       string  `database:"building_type" json:"buildingType"`
	BuildingLevel      int     `database:"building_level" json:"buildingLevel"`
	CityId             string  `database:"city_id" json:"cityId"`
	CityRow            int     `database:"city_row" json:"cityRow"`
	CityColumn         int     `database:"city_column" json:"cityColumn"`
	BuildingProduction float64 `database:"building_production" json:"buildingProduction"`
	HappinessChange    float64 `database:"happiness_change" json:"happinessChange"`
	StartTime          string  `database:"start_time" json:"startTime"`
	EndTime            string  `database:"end_time" json:"endTime"`
	// BuildingName       string  `database:"building_name" json:"buildingName"`
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
	r.HandleFunc("/cities/{session_id}/updateName", updateName).Methods("POST")
}

func getCity(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	sessionId := vars["session_id"]

	cityName := request.URL.Query()["cityName"]

	var city City

	defer func() {
		json.NewEncoder(response).Encode(city)
	}()

	var result []City

	if len(cityName) > 0 {
		database.Query(fmt.Sprintf("SELECT username, balance, population, population_capacity, city_name FROM Cities JOIN Accounts ON city_owner=player_id WHERE city_name='%s'", cityName[0]), &result)
	} else {
		database.Query(fmt.Sprintf("SELECT username, balance, population, population_capacity, city_name FROM Cities JOIN Sessions NATURAL JOIN Accounts ON city_owner=player_id WHERE session_id='%s' AND town=0", sessionId), &result)
	}

	if len(result) > 0 {
		result[0].PlayerBalance = math.Round(result[0].PlayerBalance*100) / 100
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

	cityName := request.URL.Query()["cityName"]

	var query string
	if len(cityName) > 0 {
		query = fmt.Sprintf("SELECT building_type, building_level, city_row, city_column FROM Buildings NATURAL JOIN Cities JOIN Accounts ON city_owner=player_id WHERE city_name='%s';", cityName[0])

	} else {
		query = fmt.Sprintf("SELECT building_type, building_level, city_row, city_column FROM Buildings NATURAL JOIN Cities JOIN Sessions ON city_owner=player_id WHERE session_id='%s';", sessionId)
	}

	var buildingResult []Building
	database.Query(query, &buildingResult)

	var isOwner bool
	if len(cityName) > 0 {
		database.QueryValue(
			fmt.Sprintf(
				"SELECT player_id=(SELECT player_id FROM Sessions WHERE session_id='%s') FROM Cities JOIN Sessions ON city_owner=player_id WHERE city_name='%s'", sessionId, cityName[0]), &isOwner)
	} else {
		isOwner = true
	}

	buildings = Buildings{IsOwner: isOwner, Buildings: buildingResult}
}

func getBuilding(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	sessionId := vars["session_id"]
	cityRow, _ := strconv.Atoi(vars["city_row"])
	cityColumn, _ := strconv.Atoi(vars["city_column"])

	cityName := request.URL.Query()["cityName"]

	var building []Building
	defer func() {
		if len(building) == 0 {
			json.NewEncoder(response).Encode(Building{})
		} else {
			json.NewEncoder(response).Encode(building[0])
		}
	}()

	var query string
	if len(cityName) > 0 {
		query = fmt.Sprintf("SELECT building_type, building_level, building_production, happiness_change, start_time, end_time FROM Buildings NATURAL JOIN Building_Info LEFT JOIN Builds ON Buildings.city_id=Builds.city_id AND Buildings.city_row=Builds.city_row AND Buildings.city_column=Builds.city_column WHERE Buildings.city_id=(SELECT city_id FROM Cities WHERE city_name='%s') AND Buildings.city_row=%d AND Buildings.city_column=%d;", cityName[0], cityRow, cityColumn)
	} else {
		query = fmt.Sprintf("SELECT building_type, building_level, building_production, happiness_change, start_time, end_time FROM Buildings NATURAL JOIN Building_Info LEFT JOIN Builds ON Buildings.city_id=Builds.city_id AND Buildings.city_row=Builds.city_row AND Buildings.city_column=Builds.city_column WHERE Buildings.city_id=(SELECT city_id FROM Sessions JOIN Cities ON player_id=city_owner WHERE session_id='%s' AND town=0) AND Buildings.city_row=%d AND Buildings.city_column=%d;", sessionId, cityRow, cityColumn)
	}

	database.Query(query, &building)
}

func createBuilding(response http.ResponseWriter, request *http.Request) {
	fmt.Println("Received request to /cities/createBuilding")
	vars := mux.Vars(request)
	sessionId := vars["session_id"]
	cityName := request.URL.Query()["cityName"]
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

	var query string
	if len(cityName) > 0 {
		query = fmt.Sprintf("INSERT INTO Buildings SELECT '%s', 1, city_id, %d, %d FROM Sessions JOIN Cities ON player_id=city_owner WHERE session_id='%s' AND city_name='%s'", building.BuildingType, building.CityRow, building.CityColumn, sessionId, cityName[0])
	} else {
		query = fmt.Sprintf("INSERT INTO Buildings SELECT '%s', 1, city_id, %d, %d FROM Sessions JOIN Cities ON player_id=city_owner WHERE session_id='%s' AND town=0", building.BuildingType, building.CityRow, building.CityColumn, sessionId)
	}

	result, err = database.Execute(query)
	if err != nil {
		result, err = database.Execute(
			fmt.Sprintf(
				"UPDATE Accounts SET balance = balance + (SELECT build_cost FROM Building_Info WHERE building_type='%s' AND building_level=1) WHERE player_id=(SELECT player_id FROM Sessions WHERE session_id='%s')", building.BuildingType, sessionId))
		return
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return
	}

	status = true
}

func upgradeBuilding(response http.ResponseWriter, request *http.Request) {
	fmt.Println("Received request to /cities/upgradeBuilding")
	vars := mux.Vars(request)
	sessionId := vars["session_id"]
	cityName := request.URL.Query()["cityName"]
	status := false

	defer func() {
		json.NewEncoder(response).Encode(Status{Status: status})
	}()

	var building Building
	err := json.NewDecoder(request.Body).Decode(&building)

	if err != nil {
		return
	}

	var query string
	if len(cityName) > 0 {
		query = fmt.Sprintf("UPDATE Accounts SET balance = balance - (SELECT build_cost FROM Building_Info WHERE building_type=(SELECT building_type FROM Buildings NATURAL JOIN Cities WHERE city_name='%s' AND city_row=%d AND city_column=%d) AND building_level=(SELECT building_level+1 FROM Buildings NATURAL JOIN Cities WHERE city_name='%s' AND city_row=%d AND city_column=%d)) WHERE player_id=(SELECT player_id FROM Sessions WHERE session_id='%s');", cityName[0], building.CityRow, building.CityColumn, cityName[0], building.CityRow, building.CityColumn, sessionId)
	} else {
		query = fmt.Sprintf("UPDATE Accounts SET balance = balance - (SELECT build_cost FROM Building_Info WHERE building_type=(SELECT building_type FROM Buildings NATURAL JOIN Cities JOIN Sessions ON city_owner=player_id WHERE session_id='%s' AND town=0 AND city_row=%d AND city_column=%d) AND building_level=(SELECT building_level+1 FROM Buildings NATURAL JOIN Cities JOIN Sessions ON player_id=city_owner WHERE session_id='%s' AND town=0 AND city_row=%d AND city_column=%d)) WHERE player_id=(SELECT player_id FROM Sessions WHERE session_id='%s');", sessionId, building.CityRow, building.CityColumn, sessionId, building.CityRow, building.CityColumn, sessionId)
	}
	result, err := database.Execute(query)

	if err != nil {
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return
	}

	if len(cityName) > 0 {
		query = fmt.Sprintf("UPDATE Buildings SET building_level=building_level+1 WHERE city_id=(SELECT city_id FROM Sessions JOIN Cities ON player_id=city_owner WHERE session_id='%s' AND city_name='%s') AND city_row=%d AND city_column=%d", sessionId, cityName[0], building.CityRow, building.CityColumn)
	} else {
		query = fmt.Sprintf("UPDATE Buildings SET building_level=building_level+1 WHERE city_id=(SELECT city_id FROM Sessions JOIN Cities ON player_id=city_owner WHERE session_id='%s' AND town=0) AND city_row=%d AND city_column=%d", sessionId, building.CityRow, building.CityColumn)
	}

	result, err = database.Execute(query)

	if err != nil {
		if len(cityName) > 0 {
			query = fmt.Sprintf("UPDATE Accounts SET balance = balance + (SELECT build_cost FROM Building_Info WHERE building_type=(SELECT building_type FROM Buildings NATURAL JOIN Cities WHERE city_name='%s' AND city_row=%d AND city_column=%d) AND building_level=(SELECT building_level+1 FROM Buildings NATURAL JOIN Cities WHERE city_name='%s' AND city_row=%d AND city_column=%d)) WHERE player_id=(SELECT player_id FROM Sessions WHERE session_id='%s');", cityName[0], building.CityRow, building.CityColumn, cityName[0], building.CityRow, building.CityColumn, sessionId)
		} else {
			query = fmt.Sprintf("UPDATE Accounts SET balance = balance + (SELECT build_cost FROM Building_Info WHERE building_type=(SELECT building_type FROM Buildings NATURAL JOIN Cities JOIN Sessions ON city_owner=player_id WHERE session_id='%s' AND town=0 AND city_row=%d AND city_column=%d) AND building_level=(SELECT building_level+1 FROM Buildings NATURAL JOIN Cities JOIN Sessions ON player_id=city_owner WHERE session_id='%s' AND town=0 AND city_row=%d AND city_column=%d)) WHERE player_id=(SELECT player_id FROM Sessions WHERE session_id='%s');", sessionId, building.CityRow, building.CityColumn, sessionId, building.CityRow, building.CityColumn, sessionId)
		}
		result, err = database.Execute(query)

		return
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		if len(cityName) > 0 {
			query = fmt.Sprintf("UPDATE Accounts SET balance = balance + (SELECT build_cost FROM Building_Info WHERE building_type=(SELECT building_type FROM Buildings NATURAL JOIN Cities WHERE city_name='%s' AND city_row=%d AND city_column=%d) AND building_level=(SELECT building_level+1 FROM Buildings NATURAL JOIN Cities WHERE city_name='%s' AND city_row=%d AND city_column=%d)) WHERE player_id=(SELECT player_id FROM Sessions WHERE session_id='%s');", cityName[0], building.CityRow, building.CityColumn, cityName[0], building.CityRow, building.CityColumn, sessionId)
		} else {
			query = fmt.Sprintf("UPDATE Accounts SET balance = balance + (SELECT build_cost FROM Building_Info WHERE building_type=(SELECT building_type FROM Buildings NATURAL JOIN Cities JOIN Sessions ON city_owner=player_id WHERE session_id='%s' AND town=0 AND city_row=%d AND city_column=%d) AND building_level=(SELECT building_level+1 FROM Buildings NATURAL JOIN Cities JOIN Sessions ON player_id=city_owner WHERE session_id='%s' AND town=0 AND city_row=%d AND city_column=%d)) WHERE player_id=(SELECT player_id FROM Sessions WHERE session_id='%s');", sessionId, building.CityRow, building.CityColumn, sessionId, building.CityRow, building.CityColumn, sessionId)
		}
		result, err = database.Execute(query)

		return
	}

	status = true
}

func updateName(response http.ResponseWriter, request *http.Request) {
	fmt.Println("Received request to /cities/updateName")
	vars := mux.Vars(request)
	sessionId := vars["session_id"]
	status := false

	defer func() {
		json.NewEncoder(response).Encode(Status{Status: status})
	}()

	var city CityNameChange
	err := json.NewDecoder(request.Body).Decode(&city)

	if err != nil {
		return
	}

	var query string

	if city.CityNameOriginal != "" {
		query = fmt.Sprintf("UPDATE Cities SET city_name='%s' WHERE city_owner=(SELECT player_id FROM Sessions NATURAL JOIN Accounts WHERE session_id='%s') AND city_name='%s';", city.CityNameNew, sessionId, city.CityNameOriginal)
	} else {
		query = fmt.Sprintf("UPDATE Cities SET city_name='%s' WHERE city_owner=(SELECT player_id FROM Sessions NATURAL JOIN Accounts WHERE session_id='%s') AND town=0;", city.CityNameNew, sessionId)
	}

	result, err := database.Execute(query)

	if err != nil {
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return
	}

	status = true
}
