package game

import (
	"api/database"
	"encoding/json"

	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type City struct {
	CityId     string `database:"city_id" json:"cityId"`
	Population int    `database:"population" json:"population"`
	CityName   string `database:"city_name" json:"cityName"`
	CityOwner  string `database:"city_owner" json:"cityOwner"`
}

type Buildings struct {
	IsOwner   bool        `json:"isOwner"`
	Buildings interface{} `json:"buildings"`
}

type Building struct {
	BuildingType       string  `database:"building_type"`
	BuildingLevel      string  `database:building_level`
	BuildingName       string  `database:building_name`
	CityId             string  `database:"city_id"`
	CityRow            int     `database:"city_row"`
	CityColumn         int     `database:"city_column"`
	BuildingProduction float64 `database:"building_production"`
	HappinessChange    float64 `database:"happiness_change"`
	BuildCost          float64 `database:"build_cost"`
	BuildTime          int     `database:"build_time"`
}

func HandleCityRoutes(r *mux.Router) {
	r.HandleFunc("/city/{city_id}", getCity).Methods("GET")
	r.HandleFunc("/city/{city_id}/buildings", getBuildings).Methods("GET")

	r.HandleFunc("/city/{city_id}/createBuilding", createBuilding).Methods("POST")
	r.HandleFunc("/city/{city_id}/upgradeBuilding", upgradeBuilding).Methods("POST")
}

func getCity(response http.ResponseWriter, request *http.Request) {
	log.Println("Received request to /city/{city_id}")

	vars := mux.Vars(request)
	cityId := vars["city_id"]

	var city City

	defer func() {
		json.NewEncoder(response).Encode(city)
	}()

	var result []City
	database.Query(fmt.Sprintf("SELECT * FROM Cities WHERE city_id='%s'", cityId), &result)

	city = result[0]
}

func getBuildings(response http.ResponseWriter, request *http.Request) {
	log.Println("Received request to /city/{city_id}")

	vars := mux.Vars(request)
	cityId := vars["city_id"]

	var buildings interface{}
	defer func() {
		json.NewEncoder(response).Encode(buildings)
	}()

	param := request.URL.Query()["sessionId"]
	if len(param) < 1 {
		return
	}
	sessionId := param[0]

	var buildingResult []Building
	database.Query(fmt.Sprintf("SELECT * FROM Buildings NATURAL JOIN Building_Info WHERE city_id='%s'", cityId), &buildingResult)

	var isOwner bool
	database.QueryValue(fmt.Sprintf("SELECT city_owner='%s' FROM Cities WHERE city_id='%s'", sessionId, cityId), &isOwner)

	buildings = Buildings{IsOwner: isOwner, Buildings: buildingResult}
}

func createBuilding(response http.ResponseWriter, request *http.Request) {

}

func upgradeBuilding(response http.ResponseWriter, request *http.Request) {

}
