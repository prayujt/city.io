package tests

import (
	"api/database"
	"api/game"
	"time"

	"encoding/json"
	"testing"
)

func TestArmyTrain(t *testing.T) {
	train := game.Train{
		TroopCount: 1,
	}

	response := Post("/armies/train", train)
	var result game.Status
	json.Unmarshal(response, &result)

	time.Sleep(time.Second * 1)

	if !result.Status {
		t.Error("Expected to succeed in training army")
	}

	response = Get("/armies/training")

	var trainingResult game.Training
	json.Unmarshal(response, &trainingResult)

	if trainingResult.ArmySize != 1 {
		t.Error("Expected army size of 1")
	}

	database.Execute("DELETE FROM Training")
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

	if !result1.Status {
		t.Error("Expected to succeed in training army")
	}

	response1 = Get("/armies/training")

	var trainingResult game.Training
	json.Unmarshal(response1, &trainingResult)

	if trainingResult.ArmySize != 1 {
		t.Error("Expected army size of 1")
	}

	database.Execute("DELETE FROM Training")

	march := game.March{
		ArmySize: 1,
		FromCity: "monkee monkee",
		ToCity:   "monkee",
	}

	response2 := Post("/armies/move", march)

	var result2 game.Status
	json.Unmarshal(response2, &result2)

	if !result2.Status {
		t.Error("Expected to succeed in starting march")
	}

	response2 = Get("/armies/marches")
	var marchResult []game.March
	json.Unmarshal(response2, &marchResult)

	if marchResult[0].IsAttack != true {
		t.Error("Expected an attack")
	}
}
