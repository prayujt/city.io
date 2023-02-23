package tests

import (
	"api/game"
	"time"

	"encoding/json"
	"fmt"
	"testing"
)

func TestCityGet(t *testing.T) {
	response := Get(fmt.Sprintf("/cities/%s", sessionId))

	var result game.City
	json.Unmarshal(response, &result)

	if result.CityId == "" || result.Population == 0 || result.CityName == "" {
		t.Error("Expected city to not be in null state")
	}
}

func TestCityGetFail(t *testing.T) {
	response := Get("/cities/abcdefghijklmnop")

	var result game.City
	json.Unmarshal(response, &result)

	if !(result.CityId == "" || result.Population == 0 || result.CityName == "") {
		t.Error("Expected non-existent city to be in null state")
	}
}

func TestBuildingOwnedGet(t *testing.T) {
	response := Get(fmt.Sprintf("/cities/%s/buildings", sessionId))

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
	response := Get(fmt.Sprintf("/cities/%s/buildings?username=player1", sessionId))

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
	building := game.Building{
		BuildingName:  "Hospital",
		BuildingType:  "Hospital",
		BuildingLevel: 1,
		CityRow:       0,
		CityColumn:    0,
	}

	response := Post(fmt.Sprintf("/cities/%s/createBuilding", sessionId), building)
	var result game.Status
	json.Unmarshal(response, &result)

	if !result.Status {
		t.Error("Expected to succeed in creating building")
	}

	response = Get(fmt.Sprintf("/cities/%s/buildings", sessionId))

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
	building := game.Building{
		BuildingName:  "School",
		BuildingType:  "School",
		BuildingLevel: 1,
		CityRow:       0,
		CityColumn:    0,
	}

	response := Post(fmt.Sprintf("/cities/%s/createBuilding", sessionId), building)
	var result game.Status
	json.Unmarshal(response, &result)

	if result.Status {
		t.Error("Expected to fail in creating duplicate building")
	}

	response = Get(fmt.Sprintf("/cities/%s/buildings", sessionId))

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
		BuildingName:  "Test",
		BuildingType:  "Test",
		BuildingLevel: 1,
		CityRow:       1,
		CityColumn:    1,
	}

	response := Post(fmt.Sprintf("/cities/%s/createBuilding", sessionId), building)
	var result game.Status
	json.Unmarshal(response, &result)

	if !result.Status {
		t.Error("Expected to succeed in creating building")
	}

	time.Sleep(time.Second * 2)

	response = Post(fmt.Sprintf("/cities/%s/upgradeBuilding", sessionId), building)
	var result2 game.Status
	json.Unmarshal(response, &result2)

	if !result2.Status {
		t.Error("Expected to pass in upgrading")
	}

	response = Get(fmt.Sprintf("/cities/%s/buildings", sessionId))

	var buildingUpgrade game.Buildings
	json.Unmarshal(response, &buildingUpgrade)

	if !buildingUpgrade.IsOwner {
		t.Error("Expected city to be owned by player")
	}

	var foundBuilding bool
	for _, building := range buildingUpgrade.Buildings {
		if building.CityRow == 1 && building.CityColumn == 1 && building.BuildingLevel == 2 {
			foundBuilding = true
		}
	}

	if foundBuilding == false {
		t.Error("Expected the building to upgrade, instead stayed the same")
	}

}
