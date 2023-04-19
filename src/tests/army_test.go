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

func TestMarch(t *testing.T) {
	response := Get("/armies/marches")
	var result []game.March
	json.Unmarshal(response, &result)

	if len(result) != 0 {
		t.Error("Expected no marches")
	}
}

func TestStartMarchAttackPass(t *testing.T) {
	acc5 := login.Account{
		Username: "rawr5",
		Password: "rawr4",
	}

	Post("/login/createAccount", acc5)

	response := Post("/login/createSession", acc5)

	var session3 login.JWT
	json.Unmarshal(response, &session3)

	if session3.Token == "" {
		fmt.Println("Failed to initialize token for tests")
		return
	}

	city := game.CityNameChange{
		CityNameNew: "monkee3",
	}

	response1 := Post("/cities/updateName", city, session3.Token)
	var result game.Status

	json.Unmarshal(response1, &result)

	if !result.Status {
		t.Error("Expected to succeed in changing name")
	}

	building := game.Building{
		BuildingType:  "Barracks",
		BuildingLevel: 1,
		CityRow:       1,
		CityColumn:    1,
	}

	response2 := Post("/cities/createBuilding", building, session3.Token)
	var result2 game.Status
	json.Unmarshal(response2, &result2)

	if !result2.Status {
		t.Error("Expected to succeed in creating building")
	}

	train := game.Train{
		TroopCount: 1,
	}

	response3 := Post("/armies/train", train, session3.Token)
	var result3 game.Status
	json.Unmarshal(response3, &result3)

	if !result3.Status {
		t.Error("Expected to succeed")
	}

	response = Get("/armies/training", session3.Token)

	var trainingResult game.Training
	json.Unmarshal(response, &trainingResult)

	database.Execute("DELETE FROM Training")

	if trainingResult.ArmySize != 1 {
		t.Error("Expected army size of 1")
	}

	march := game.March{
		ArmySize: 1,
		FromCity: "monkee3",
		ToCity:   "monkee monkee",
	}

	response4 := Post("/armies/move", march, session3.Token)
	var result4 game.Status
	json.Unmarshal(response4, &result4)

	if !result4.Status {
		t.Error("Expected to succeed in starting march")
	}

	response4 = Get("/armies/marches", session3.Token)
	var marchResult []game.March
	json.Unmarshal(response4, &marchResult)

	if marchResult[0].IsAttack != true {
		t.Error("Expected an attack")
	}
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
		CityRow:       1,
		CityColumn:    1,
	}

	response2 := Post("/cities/createBuilding", building, session3.Token)
	var result2 game.Status
	json.Unmarshal(response2, &result2)

	if !result2.Status {
		t.Error("Expected to succeed in creating building")
	}

	train := game.Train{
		TroopCount: 1,
	}

	response3 := Post("/armies/train", train, session3.Token)
	var result3 game.Status
	json.Unmarshal(response3, &result3)

	if !result3.Status {
		t.Error("Expected to succeed")
	}

	response = Get("/armies/training", session3.Token)

	var trainingResult game.Training
	json.Unmarshal(response, &trainingResult)

	if trainingResult.ArmySize != 1 {
		t.Error("Expected army size of 1")
	}
}

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
