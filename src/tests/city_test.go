package tests

import (
	"api/game"
	"api/login"
	"strings"
	"time"

	"encoding/json"
	"fmt"
	"testing"
)

func TestCityGet(t *testing.T) {
	response := Get(fmt.Sprintf("/cities/stats"))

	var result game.City
	json.Unmarshal(response, &result)

	if result.Population == 0 || result.CityName == "" {
		t.Error("Expected city to not be in null state")
	}
}

func TestCityGetFail(t *testing.T) {
	response := Get("/cities/stats", "abcdefghijklmnop")

	var result game.City
	json.Unmarshal(response, &result)

	if !(result.Population == 0 || result.CityName == "") {
		t.Error("Expected non-existent city to be in null state")
	}
}

func TestBuildingOwnedGet(t *testing.T) {
	response := Get("/cities/buildings")

	var result game.Buildings
	json.Unmarshal(response, &result)

	if !result.IsOwner {
		t.Error("Expected city to be owned by player")
	}

	if len(result.Buildings) == 0 {
		t.Error("Expected a building to exist in the city")
	}
}

func TestBuildingNotOwnedGet(t *testing.T) {
	response := Get("/cities")

	var cities []game.Ownership
	json.Unmarshal(response, &cities)

	var unownedCity string
	for _, pair := range cities {
		if pair.CityOwner == "player1" {
			unownedCity = pair.CityName
		}
	}

	response = Get(fmt.Sprintf("/cities/buildings?cityName=%s", strings.ReplaceAll(unownedCity, " ", "+")))

	var result game.Buildings
	json.Unmarshal(response, &result)

	if result.IsOwner {
		t.Error("Expected city to not be owned by player")
	}

	if len(result.Buildings) == 0 {
		t.Error("Expected a building to exist in the city")
	}
}

func TestBuildingCreate(t *testing.T) {
	acc6 := login.Account{
		Username: "rawr6",
		Password: "rawr4",
	}

	Post("/login/createAccount", acc6)

	response := Post("/login/createSession", acc6)

	var session3 login.JWT
	json.Unmarshal(response, &session3)

	if session3.Token == "" {
		fmt.Println("Failed to initialize token for tests")
		return
	}

	building := game.Building{
		BuildingType:  "Hospital",
		BuildingLevel: 1,
		CityRow:       0,
		CityColumn:    0,
	}

	response2 := Post("/cities/createBuilding", building, session3.Token)
	var result game.Status
	json.Unmarshal(response2, &result)

	if !result.Status {
		t.Error("Expected to succeed in creating building")
	}

	response = Get("/cities/buildings", session3.Token)

	var buildingsResult game.Buildings
	json.Unmarshal(response, &buildingsResult)

	if !buildingsResult.IsOwner {
		t.Error("Expected city to be owned by player")
	}

	if len(buildingsResult.Buildings) != 2 {
		t.Error("Expected city to have exactly two buildings")
	}
}

func TestBuildingCreateDuplicate(t *testing.T) {
	acc := login.Account{
		Username: "rawr7",
		Password: "rawr4",
	}

	Post("/login/createAccount", acc)

	response := Post("/login/createSession", acc)

	var session login.JWT
	json.Unmarshal(response, &session)

	if session.Token == "" {
		fmt.Println("Failed to initialize token for tests")
		return
	}

	building1 := game.Building{
		BuildingType:  "Hospital",
		BuildingLevel: 1,
		CityRow:       0,
		CityColumn:    0,
	}

	response2 := Post("/cities/createBuilding", building1, session.Token)
	var result1 game.Status
	json.Unmarshal(response2, &result1)

	if !result1.Status {
		t.Error("Expected to succeed in creating building")
	}

	response = Get("/cities/buildings", session.Token)

	var buildingsResult1 game.Buildings
	json.Unmarshal(response, &buildingsResult1)

	if !buildingsResult1.IsOwner {
		t.Error("Expected city to be owned by player")
	}

	if len(buildingsResult1.Buildings) != 2 {
		t.Error("Expected city to have exactly two buildings")
	}

	building := game.Building{
		BuildingType:  "School",
		BuildingLevel: 1,
		CityRow:       0,
		CityColumn:    0,
	}

	response3 := Post("/cities/createBuilding", building, session.Token)
	var result game.Status
	json.Unmarshal(response3, &result)

	if result.Status {
		t.Error("Expected to fail in creating duplicate building")
	}

	response = Get("/cities/buildings", session.Token)

	var buildingsResult game.Buildings
	json.Unmarshal(response, &buildingsResult)

	if !buildingsResult.IsOwner {
		t.Error("Expected city to be owned by player")
	}

	for _, building := range buildingsResult.Buildings {
		if building.BuildingType == "School" {
			t.Error("Expected the building to remain the same, instead got new building")
		}
	}
}

func TestUpgradeBuilding(t *testing.T) {
	building := game.Building{
		BuildingType:  "Test",
		BuildingLevel: 1,
		CityRow:       0,
		CityColumn:    1,
	}

	response := Post("/cities/createBuilding", building)
	var result game.Status
	json.Unmarshal(response, &result)

	if !result.Status {
		t.Error("Expected to succeed in creating building")
	}

	time.Sleep(time.Second * 2)

	response = Post("/cities/upgradeBuilding", building)
	var result2 game.Status
	json.Unmarshal(response, &result2)

	if !result2.Status {
		t.Error("Expected to pass in upgrading")
	}

	response = Get("/cities/buildings")

	var buildingUpgrade game.Buildings
	json.Unmarshal(response, &buildingUpgrade)

	if !buildingUpgrade.IsOwner {
		t.Error("Expected city to be owned by player")
	}

	var foundBuilding bool
	for _, building2 := range buildingUpgrade.Buildings {
		if building2.CityRow == building.CityRow && building2.CityColumn == building.CityColumn && building2.BuildingLevel == building.BuildingLevel+1 {
			foundBuilding = true
		}
	}

	if foundBuilding == false {
		t.Error("Expected the building to upgrade, instead stayed the same")
	}

}

func TestNameChange(t *testing.T) {
	acc3 := login.Account{
		Username: "rawr",
		Password: "rawr",
	}

	Post("/login/createAccount", acc3)

	response := Post("/login/createSession", acc3)

	var session3 login.JWT
	json.Unmarshal(response, &session3)

	if session3.Token == "" {
		fmt.Println("Failed to initialize token for tests")
		return
	}

	city := game.CityNameChange{
		CityNameNew: "monkee5",
	}

	response1 := Post("/cities/updateName", city, session3.Token)
	var result game.Status

	json.Unmarshal(response1, &result)

	if !result.Status {
		t.Error("Expected to succeed in changing name")
	}
}
