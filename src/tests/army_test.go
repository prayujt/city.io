package tests

import (
	"api/database"
	"api/game"
	"api/login"
	"time"

	"encoding/json"
	"fmt"
	"testing"
)

func TestArmyTrainFail(t *testing.T) {
	train := game.Train{
		TroopCount: 1,
	}

	response := Post("/armies/train", train)
	var result game.Status
	json.Unmarshal(response, &result)

	time.Sleep(time.Second * 1)

	if result.Status {
		t.Error("Expected not to succeed")
	}
}

func TestMarch(t *testing.T) {
	response := Get("/armies/marches")
	var result []game.March
	json.Unmarshal(response, &result)

	if len(result) != 0 {
		t.Error("Expected no marches")
	}
}

func TestStartMarchAttack(t *testing.T) {
	train := game.Train{
		TroopCount: 1,
	}

	response1 := Post("/armies/train", train)
	var result1 game.Status
	json.Unmarshal(response1, &result1)

	if result1.Status {
		t.Error("Expected to succeed in training army")
	}

	response1 = Get("/armies/training")

	var trainingResult game.Training
	json.Unmarshal(response1, &trainingResult)

	if trainingResult.ArmySize != 0 {
		t.Error("Expected army size of 1")
	}

	database.Execute("DELETE FROM Training")
}

func TestArmyTrainSuccess(t *testing.T) {
	acc4 := login.Account{
		Username: "rawr4",
		Password: "rawr4",
	}

	Post("/login/createAccount", acc4)

	response := Post("/login/createSession", acc4)

	var session3 login.JWT
	json.Unmarshal(response, &session3)

	if session3.Token == "" {
		fmt.Println("Failed to initialize token for tests")
		return
	}

	building := game.Building{
		BuildingType:  "Barracks",
		BuildingLevel: 1,
		CityRow:       0,
		CityColumn:    0,
	}

	response2 := Post("/cities/createBuilding", building)
	var result2 game.Status
	json.Unmarshal(response2, &result2)

	if !result2.Status {
		t.Error("Expected to succeed in creating building")
	}

	train := game.Train{
		TroopCount: 1,
	}

	response3 := Post("/armies/train", train)
	var result3 game.Status
	json.Unmarshal(response3, &result3)

	if !result3.Status {
		t.Error("Expected to succeed")
	}

	response = Get("/armies/training")

	var trainingResult game.Training
	json.Unmarshal(response, &trainingResult)

	if trainingResult.ArmySize != 1 {
		t.Error("Expected army size of 1")
	}
}
