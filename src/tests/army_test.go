package tests

import (
	"api/game"
	"time"

	"encoding/json"
	"testing"
)

func TestArmyTrain(t *testing.T) {
	train := game.Train{
		TroopCount: 10,
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

	if trainingResult.ArmySize != 10 {
		t.Error("Expected army size of 100")
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

func TestStartMarch(t *testing.T) {

	march := game.March{
		ArmySize: 10,
		FromCity: "monkee monkee",
		ToCity:   "monkee monkee",
	}

	response := Post("/armies/move", march)
	var result game.Status
	json.Unmarshal(response, &result)

	if result.Status {
		t.Error("Expected to succeed in starting march")
	}

	time.Sleep(time.Second * 2)

	response = Get("/armies/marches")
	var marchResult game.March
	json.Unmarshal(response, &marchResult)

	if marchResult.IsAttack == true {
		t.Error("Expected a move")
	}
}
