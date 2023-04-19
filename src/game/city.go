package game

import (
	"api/auth"
	"api/database"
	"encoding/json"
	"math"

	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type City struct {
	CityName           string  `database:"city_name" json:"cityName"`
	Population         int     `database:"population" json:"population"`
	PopulationCapacity int     `database:"population_capacity" json:"populationCapacity"`
	PlayerBalance      float64 `database:"balance" json:"playerBalance"`
	CityOwner          string  `database:"username" json:"cityOwner"`
	ArmySize           int     `database:"army_size" json:"armySize"`
	Happiness          int     `database:"happiness_total" json:"happinessTotal"`
}

type CityStats struct {
	CityName       string `database:"city_name" json:"cityName"`
	ArmySize       int    `database:"army_size" json:"armySize"`
	CityProduction int    `database:"city_production" json:"cityProduction"`
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
	HappinessChange    int     `database:"happiness_change" json:"happinessChange"`
	StartTime          string  `database:"start_time" json:"startTime"`
	EndTime            string  `database:"end_time" json:"endTime"`
	UpgradeCost        float64 `json:"upgradeCost"`
	UpgradedProduction float64 `json:"upgradedProduction"`
	UpgradeTime        int     `json:"upgradeTime"`
	UpgradedHappiness  int     `json:"upgradeHappniess"`
}

type NewBuilding struct {
	BuildCost          float64 `database:"build_cost" json:"buildCost"`
	BuildTime          int     `database:"build_time" json:"buildTime"`
	BuildingType       string  `database:"building_type" json:"buildingType"`
	BuildingProduction float64 `database:"building_production" json:"buildingProduction"`
	HappinessChange    int     `database:"happiness_change" json:"happinessChange"`
}

type CityInfo struct {
	CityName        string `database:"city_name" json:"cityName"`
	ProductionTotal int    `database:"total_production" json:"totalProduction"`
	ArmySize        int    `database:"army_size" json:"armySize"`
	Population      int    `database:"population_total" json:"totalPopulation"`
}

type Status struct {
	Status bool `json:"status"`
}

func HandleCityRoutes(r *mux.Router) {
	r.HandleFunc("/cities/buildings/available", getAllBuildings).Methods("GET")
	r.HandleFunc("/cities/stats", getCityStats).Methods("GET")
	r.HandleFunc("/cities/territory", getTerritory).Methods("GET")
	r.HandleFunc("/cities/buildings", getBuildings).Methods("GET")
	r.HandleFunc("/cities/buildings/{city_row}/{city_column}", getBuilding).Methods("GET")
	r.HandleFunc("/cities/production", getProduction).Methods("GET")

	r.HandleFunc("/cities/createBuilding", createBuilding).Methods("POST")
	r.HandleFunc("/cities/upgradeBuilding", upgradeBuilding).Methods("POST")
	r.HandleFunc("/cities/updateName", updateName).Methods("POST")
	r.HandleFunc("/cities/destroyBuilding", destroyBuilding).Methods("Post")
}

func getAllBuildings(response http.ResponseWriter, request *http.Request) {
	var buildings []NewBuilding

	defer func() {
		json.NewEncoder(response).Encode(buildings)
	}()

	database.Query(
		`
		SELECT building_type, build_cost, build_time, building_production, happiness_change
		FROM Building_Info
		WHERE building_level=1 AND building_type != 'Test' AND building_type != 'City Hall'
		`,
		&buildings)
}

func getCityStats(response http.ResponseWriter, request *http.Request) {
	var city City

	defer func() {
		json.NewEncoder(response).Encode(city)
	}()

	if request.Header["Token"] == nil {
		return
	}

	claims, err := auth.ParseJWT(request.Header["Token"][0])

	if err != nil {
		return
	}

	cityName := request.URL.Query()["cityName"]

	var result []City

	if len(cityName) > 0 {
		database.Query(
			fmt.Sprintf(
				`
				SELECT username, balance, population, 
				(SELECT SUM(happiness_change) 
				FROM Buildings JOIN Building_Info ON Buildings.building_type=Building_Info.building_type AND Buildings.building_level=Building_Info.building_level
				WHERE city_id =(SELECT city_id FROM Cities where city_name='%s')) AS happiness_total,
				(SELECT SUM(population_capacity_change) 
				FROM Buildings JOIN Building_Info ON Buildings.building_type=Building_Info.building_type AND Buildings.building_level=Building_Info.building_level
				WHERE city_id =(SELECT city_id FROM Cities where city_name='%s')) AS population_capacity,
				IF(username = '%s', army_size, -1) AS army_size, city_name
				FROM Cities JOIN Accounts ON city_owner=player_id
				WHERE city_name='%s'
				`,
				cityName[0], cityName[0], claims["username"], cityName[0]),
			&result)
	} else {
		database.Query(
			fmt.Sprintf(
				`
				SELECT username, balance, population, population_capacity, army_size, city_name
				FROM Cities JOIN Accounts ON city_owner=player_id
				WHERE player_id='%s' AND town=0
				`,
				claims["playerId"]),
			&result)
	}

	if len(result) > 0 {
		result[0].PlayerBalance = math.Round(result[0].PlayerBalance*100) / 100
		city = result[0]
	}
}

func getProduction(response http.ResponseWriter, request *http.Request) {
	var city []CityInfo

	defer func() {
		json.NewEncoder(response).Encode(city)
	}()

	if request.Header["Token"] == nil {
		return
	}

	claims, err := auth.ParseJWT(request.Header["Token"][0])

	if err != nil {
		return
	}

	database.Query(
		fmt.Sprintf(
			`
			SELECT city_name, SUM(building_production) as total_production, MAX(army_size) as army_size, MAX(population) as population_total
			FROM Building_Info
			JOIN Buildings ON Building_Info.building_type = Buildings.building_type AND Building_Info.building_level = Buildings.building_level
            NATURAL JOIN Cities WHERE city_owner = '%s'
            GROUP BY city_name
			`,
			claims["playerId"]),
		&city)
}

func getTerritory(response http.ResponseWriter, request *http.Request) {
	var territory []CityStats

	defer func() {
		json.NewEncoder(response).Encode(territory)
	}()

	if request.Header["Token"] == nil {
		return
	}

	claims, err := auth.ParseJWT(request.Header["Token"][0])

	if err != nil {
		return
	}

	database.Query(
		fmt.Sprintf(
			`
			SELECT city_name, any_value(army_size) AS army_size, SUM(building_production) AS city_production
			FROM Building_Ownership
			WHERE player_id='%s' GROUP BY city_name
			`,
			claims["playerId"]),
		&territory)
}

func getBuildings(response http.ResponseWriter, request *http.Request) {
	var buildings Buildings
	defer func() {
		json.NewEncoder(response).Encode(buildings)
	}()

	if request.Header["Token"] == nil {
		return
	}

	claims, err := auth.ParseJWT(request.Header["Token"][0])

	if err != nil {
		return
	}

	cityName := request.URL.Query()["cityName"]

	var query string
	if len(cityName) > 0 {
		query = fmt.Sprintf(
			`
			SELECT building_type, building_level, city_row, city_column
			FROM Buildings NATURAL JOIN Cities
			WHERE city_name='%s'
			`,
			cityName[0])
	} else {
		query = fmt.Sprintf(
			`
			SELECT building_type, building_level, city_row, city_column
			FROM Buildings NATURAL JOIN Cities
			WHERE city_owner='%s' AND town=0
			`,
			claims["playerId"])
	}

	var buildingResult []Building
	database.Query(query, &buildingResult)

	var isOwner bool
	if len(cityName) > 0 {
		database.QueryValue(
			fmt.Sprintf(
				`
				SELECT city_owner='%s'
				FROM Cities
				WHERE city_name='%s'
				`,
				claims["playerId"], cityName[0]),
			&isOwner)
	} else {
		isOwner = true
	}

	buildings = Buildings{IsOwner: isOwner, Buildings: buildingResult}
}

func getBuilding(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
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

	if request.Header["Token"] == nil {
		return
	}

	claims, err := auth.ParseJWT(request.Header["Token"][0])

	if err != nil {
		return
	}

	var query string
	if len(cityName) > 0 {
		query = fmt.Sprintf(
			`
			SELECT building_type, building_level, building_production, happiness_change, start_time, end_time
			FROM Buildings NATURAL JOIN Building_Info LEFT JOIN Builds
				ON
				Buildings.city_id=Builds.city_id
				AND
				Buildings.city_row=Builds.city_row
				AND
				Buildings.city_column=Builds.city_column
			WHERE
				Buildings.city_id=
					(
					SELECT city_id FROM Cities WHERE city_name='%s'
					)
				AND
				Buildings.city_row=%d AND Buildings.city_column=%d
			`,
			cityName[0], cityRow, cityColumn)
	} else {
		query = fmt.Sprintf(
			`
			SELECT building_type, building_level, building_production, happiness_change, start_time, end_time
			FROM Buildings NATURAL JOIN Building_Info LEFT JOIN Builds
			ON
				Buildings.city_id=Builds.city_id
				AND
				Buildings.city_row=Builds.city_row
				AND
				Buildings.city_column=Builds.city_column
			WHERE Buildings.city_id=
				(SELECT city_id FROM Cities WHERE city_owner='%s' AND town=0)
				AND
				Buildings.city_row=%d
				AND
				Buildings.city_column=%d
			`,
			claims["playerId"], cityRow, cityColumn)
	}

	database.Query(query, &building)

	var upgradeBuilding []NewBuilding

	if len(building) == 0 {
		return
	}

	if building[0].BuildingType == "City Hall" || building[0].BuildingLevel == 10 {
		return
	}

	database.Query(
		fmt.Sprintf(
			`
			SELECT build_cost, build_time, building_production, happiness_change, population_capacity_change 
			FROM Building_Info WHERE building_type='%s' AND building_level=%d
			`,
			building[0].BuildingType, building[0].BuildingLevel+1),
		&upgradeBuilding)

	if len(upgradeBuilding) == 0 {
		return
	}

	building[0].UpgradeCost = upgradeBuilding[0].BuildCost
	building[0].UpgradedProduction = upgradeBuilding[0].BuildingProduction
	building[0].UpgradedHappiness = upgradeBuilding[0].HappinessChange
	building[0].UpgradeTime = upgradeBuilding[0].BuildTime
}

func createBuilding(response http.ResponseWriter, request *http.Request) {
	cityName := request.URL.Query()["cityName"]
	status := false

	defer func() {
		json.NewEncoder(response).Encode(Status{Status: status})
	}()

	if request.Header["Token"] == nil {
		return
	}

	claims, err := auth.ParseJWT(request.Header["Token"][0])

	if err != nil {
		return
	}

	var building Building
	err = json.NewDecoder(request.Body).Decode(&building)

	if err != nil {
		return
	}

	result, err := database.Execute(
		fmt.Sprintf(
			`
			UPDATE Accounts
			SET balance = balance -
				(SELECT build_cost FROM Building_Info WHERE building_type='%s' AND building_level=1)
			WHERE player_id='%s'
			`,
			building.BuildingType, claims["playerId"]))

	if err != nil {
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return
	}

	var query string
	if len(cityName) > 0 {
		query = fmt.Sprintf(
			`
			INSERT INTO Buildings
			SELECT '%s', 1, city_id, %d, %d
			FROM Cities WHERE city_owner='%s' AND city_name='%s'
			`,
			building.BuildingType, building.CityRow, building.CityColumn, claims["playerId"], cityName[0])
	} else {
		query = fmt.Sprintf(
			`
			INSERT INTO Buildings
			SELECT '%s', 1, city_id, %d, %d
			FROM Cities WHERE city_owner='%s' AND town=0
			`,
			building.BuildingType, building.CityRow, building.CityColumn, claims["playerId"])
	}

	result, err = database.Execute(query)

	if err != nil {
		result, _ = database.Execute(
			fmt.Sprintf(
				`
				UPDATE Accounts
				SET balance = balance +
					(SELECT build_cost FROM Building_Info WHERE building_type='%s' AND building_level=1)
				WHERE player_id='%s'
				`,
				building.BuildingType, claims["playerId"]))
		return
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return
	}

	status = true
}

func destroyBuilding(response http.ResponseWriter, request *http.Request) {
	cityName := request.URL.Query()["cityName"]
	status := false

	defer func() {
		json.NewEncoder(response).Encode(Status{Status: status})
	}()

	if request.Header["Token"] == nil {
		return
	}

	claims, err := auth.ParseJWT(request.Header["Token"][0])

	if err != nil {
		log.Println(err)
		return
	}

	var building Building
	err = json.NewDecoder(request.Body).Decode(&building)

	if err != nil {
		log.Println(err)
		return
	}

	var query string
	if len(cityName) > 0 {
		query = fmt.Sprintf(
			`
			UPDATE Accounts
			SET balance = balance +
				(SELECT SUM(build_cost)/2 AS total_cost
				FROM Building_Info
				WHERE building_level <=
					(SELECT building_level FROM Buildings WHERE city_id=(SELECT city_id FROM Cities WHERE city_name='%s' AND city_owner='%s') AND city_row=%d and city_column=%d))
			WHERE player_id='%s'
			`,
			cityName[0], claims["playerId"], building.CityRow, building.CityColumn, claims["playerId"])
	} else {
		query = fmt.Sprintf(
			`
			UPDATE Accounts
			SET balance = balance +
				(SELECT SUM(build_cost)/2 AS total_cost
				FROM Building_Info
				WHERE building_level <=
					(SELECT building_level FROM Buildings WHERE city_id=(SELECT city_id FROM Cities WHERE city_owner='%s' AND town=0) AND city_row=%d and city_column=%d))
			WHERE player_id='%s'
			`,
			claims["playerId"], building.CityRow, building.CityColumn, claims["playerId"])
	}

	result, err := database.Execute(query)

	if err != nil {
		log.Println(err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		log.Println(err)
		return
	}

	if len(cityName) > 0 {
		query = fmt.Sprintf(
			`
		DELETE FROM Buildings
		WHERE city_row=%d AND city_column=%d
		AND city_id=(SELECT city_id FROM Cities WHERE city_owner='%s' AND city_name='%s')
		`,
			building.CityRow, building.CityColumn, claims["playerId"], cityName[0])
	} else {
		query = fmt.Sprintf(
			`
		DELETE FROM Buildings
		WHERE city_row=%d AND city_column=%d
		AND city_id=(SELECT city_id FROM Cities WHERE city_owner='%s' AND town=0)
		`,
			building.CityRow, building.CityColumn, claims["playerId"])
	}

	result, err = database.Execute(query)

	if err != nil {
		if len(cityName) > 0 {
			result, _ = database.Execute(
				fmt.Sprintf(
					`
					UPDATE Accounts
					SET balance = balance -
					(SELECT SUM(build_cost)/2 AS total_cost
					FROM Building_Info
					WHERE building_level <=
					(SELECT building_level FROM Buildings WHERE city_id=(SELECT city_id FROM Cities WHERE city_name='%s') AND city_row=%d and city_column=%d AND city_owner='%s'))
					WHERE player_id='%s'
					`,
					cityName[0], building.CityRow, building.CityColumn, claims["playerId"], claims["playerId"]))
		} else {
			result, _ = database.Execute(
				fmt.Sprintf(
					`
							UPDATE Accounts
							SET balance = balance -
								(SELECT SUM(build_cost)/2 AS total_cost
								FROM Building_Info
								WHERE building_level <=
									(SELECT building_level FROM Buildings WHERE city_id=(SELECT city_id FROM Cities WHERE city_owner='%s' AND town=0) AND city_row=%d and city_column=%d))
							WHERE player_id='%s'
							`,
					claims["playerId"], building.CityRow, building.CityColumn, claims["playerId"]))
		}
		return
	}

	rowsAffected, err = result.RowsAffected()

	if err != nil || rowsAffected == 0 {
		return
	}

	status = true
}

func upgradeBuilding(response http.ResponseWriter, request *http.Request) {
	cityName := request.URL.Query()["cityName"]
	status := false

	defer func() {
		json.NewEncoder(response).Encode(Status{Status: status})
	}()

	if request.Header["Token"] == nil {
		return
	}

	claims, err := auth.ParseJWT(request.Header["Token"][0])

	if err != nil {
		return
	}

	var building Building
	err = json.NewDecoder(request.Body).Decode(&building)

	if err != nil {
		return
	}

	var query string
	if len(cityName) > 0 {
		query = fmt.Sprintf(
			`
			UPDATE Accounts
			SET balance = balance -
				(SELECT build_cost FROM Building_Info
					WHERE building_type=
						(
						SELECT building_type
						FROM (SELECT * FROM Building_Ownership) AS TempTable1
						WHERE
							city_name='%s'
							AND
							city_row=%d
							AND
							city_column=%d
						)
					AND
					building_level=
						(
						SELECT building_level+1
						FROM (SELECT * FROM Building_Ownership) AS TempTable2
						WHERE
							city_name='%s'
							AND
							city_row=%d
							AND
							city_column=%d
						)
				)
			WHERE player_id='%s'
			`,
			cityName[0], building.CityRow, building.CityColumn, cityName[0], building.CityRow, building.CityColumn, claims["playerId"])
	} else {
		query = fmt.Sprintf(
			`
			UPDATE Accounts
			SET balance = balance -
				(
				SELECT build_cost FROM Building_Info
				WHERE building_type=
					(
					SELECT building_type
					FROM (SELECT * FROM Building_Ownership) AS TempTable1
					WHERE
						player_id='%s'
						AND
						town=0
						AND
						city_row=%d
						AND
						city_column=%d
					)
				AND
				building_level=
					(
					SELECT building_level+1
					FROM (SELECT * FROM Building_Ownership) AS TempTable2
					WHERE
						player_id='%s'
						AND
						town=0
						AND
						city_row=%d
						AND city_column=%d
					)
				)
			WHERE player_id='%s'
			`,
			claims["playerId"], building.CityRow, building.CityColumn, claims["playerId"], building.CityRow, building.CityColumn, claims["playerId"])
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
		query = fmt.Sprintf(
			`
			UPDATE Buildings
			SET building_level=building_level+1
			WHERE city_id=
				(
				SELECT city_id
				FROM Cities
				WHERE city_owner='%s' AND city_name='%s'
				)
			AND
			city_row=%d
			AND
			city_column=%d
			`,
			claims["playerId"], cityName[0], building.CityRow, building.CityColumn)
	} else {
		query = fmt.Sprintf(
			`
			UPDATE Buildings
			SET building_level=building_level+1
			WHERE
			city_id=
				(
				SELECT city_id FROM Cities WHERE city_owner='%s' AND town=0
				)
			AND
			city_row=%d
			AND
			city_column=%d
			`,
			claims["playerId"], building.CityRow, building.CityColumn)
	}

	result, err = database.Execute(query)

	if err != nil {
		if len(cityName) > 0 {
			query = fmt.Sprintf(
				`
				UPDATE Accounts
				SET balance = balance +
					(
						SELECT build_cost FROM Building_Info
						WHERE
							building_type=
								(
								SELECT building_type
								FROM Building_Ownership
								WHERE
									city_name='%s'
									AND
									city_row=%d
									AND
									city_column=%d
								)
							AND
							building_level=
								(
								SELECT building_level+1
								FROM Building_Ownership
								WHERE
									city_name='%s'
									AND
									city_row=%d
									AND
									city_column=%d
								)
					)
				WHERE player_id='%s'
				`,
				cityName[0], building.CityRow, building.CityColumn, cityName[0], building.CityRow, building.CityColumn, claims["playerId"])
		} else {
			query = fmt.Sprintf(
				`
				UPDATE Accounts
				SET balance = balance +
					(
					SELECT build_cost FROM Building_Info
					WHERE
						building_type=
							(
							SELECT building_type
							FROM Building_Ownership
							WHERE
								player_id='%s'
								AND
								town=0
								AND
								city_row=%d
								AND
								city_column=%d
							)
						AND
						building_level=
							(
							SELECT building_level+1
							FROM Building_Ownership
							WHERE
								player_id='%s'
								AND
								town=0
								AND
								city_row=%d
								AND city_column=%d
							)
					)
				WHERE player_id='%s'
				`,
				claims["playerId"], building.CityRow, building.CityColumn, claims["playerId"], building.CityRow, building.CityColumn, claims["playerId"])
		}

		result, err = database.Execute(query)

		return
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		if len(cityName) > 0 {
			query = fmt.Sprintf(
				`
				UPDATE Accounts
				SET balance = balance +
					(
						SELECT build_cost FROM Building_Info
						WHERE
							building_type=
								(
								SELECT building_type
								FROM Building_Ownership
								WHERE
									city_name='%s'
									AND
									city_row=%d
									AND
									city_column=%d
								)
							AND
							building_level=
								(
								SELECT building_level+1
								FROM Building_Ownership
								WHERE
									city_name='%s'
									AND
									city_row=%d
									AND
									city_column=%d
								)
					)
				WHERE player_id='%s'
				`,
				cityName[0], building.CityRow, building.CityColumn, cityName[0], building.CityRow, building.CityColumn, claims["playerId"])
		} else {
			query = fmt.Sprintf(
				`
				UPDATE Accounts
				SET balance = balance +
					(
					SELECT build_cost FROM Building_Info
					WHERE
						building_type=
							(
							SELECT building_type
							FROM Building_Ownership
							WHERE
								player_id='%s'
								AND
								town=0
								AND
								city_row=%d
								AND
								city_column=%d
							)
						AND
						building_level=
							(
							SELECT building_level+1
							FROM Building_Ownership
							WHERE
								player_id='%s'
								AND
								town=0
								AND
								city_row=%d
								AND city_column=%d
							)
					)
				WHERE player_id='%s'
				`,
				claims["playerId"], building.CityRow, building.CityColumn, claims["playerId"], building.CityRow, building.CityColumn, claims["playerId"])
		}

		result, err = database.Execute(query)
		return
	}

	status = true
}

func updateName(response http.ResponseWriter, request *http.Request) {
	status := false

	defer func() {
		json.NewEncoder(response).Encode(Status{Status: status})
	}()

	if request.Header["Token"] == nil {
		return
	}

	claims, err := auth.ParseJWT(request.Header["Token"][0])

	if err != nil {
		return
	}

	var city CityNameChange
	err = json.NewDecoder(request.Body).Decode(&city)

	if err != nil {
		return
	}

	var query string

	if city.CityNameOriginal != "" {
		query = fmt.Sprintf(
			`
			UPDATE Cities
			SET city_name='%s'
			WHERE
				city_owner='%s'
				AND
				city_name='%s'
			`,
			city.CityNameNew, claims["playerId"], city.CityNameOriginal)
	} else {
		query = fmt.Sprintf(
			`
			UPDATE Cities
			SET city_name='%s'
			WHERE
				city_owner='%s'
				AND
				town=0
			`,
			city.CityNameNew, claims["playerId"])
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
