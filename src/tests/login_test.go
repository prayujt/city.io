package tests

import (
	"api/login"

	"encoding/json"
	"testing"
)

func TestAccountCreate(t *testing.T) {
	acc := login.Account{
		Username: "test",
		Password: "test",
	}

	response := Post("/login/createAccount", acc)

	var result login.Status
	json.Unmarshal(response, &result)

	if result.Status != true {
		t.Error("Expected account creation to return true, instead got false")
	}
}

func TestDuplicateAccount(t *testing.T) {
	acc := login.Account{
		Username: "test",
		Password: "test",
	}

	response := Post("/login/createAccount", acc)

	var result login.Status
	json.Unmarshal(response, &result)

	if result.Status != false {
		t.Error("Expected duplicate account creation to return false, instead got true")
	}
}

func TestAccountLogin(t *testing.T) {
	acc := login.Account{
		Username: "test",
		Password: "test",
	}

	response := Post("/login/createSession", acc)

	var jwt login.JWT
	json.Unmarshal(response, &jwt)

	if jwt.Token == "" {
		t.Error("Expected session creation to return true, instead got false")
	}
}

func TestIncorrectUsername(t *testing.T) {
	acc := login.Account{
		Username: "test2",
		Password: "test",
	}

	response := Post("/login/createSession", acc)

	var jwt login.JWT
	json.Unmarshal(response, &jwt)

	if jwt.Token != "" {
		t.Error("Expected session creation with incorrect username to return false, instead got true")
	}
}

func TestIncorrectPassword(t *testing.T) {
	acc := login.Account{
		Username: "test",
		Password: "test2",
	}

	response := Post("/login/createSession", acc)

	var jwt login.JWT
	json.Unmarshal(response, &jwt)

	if jwt.Token != "" {
		t.Error("Expected session creation with incorrect password to return false, instead got true")
	}
}

func TestIncorrectUsernamePassword(t *testing.T) {
	acc := login.Account{
		Username: "testing",
		Password: "testing",
	}

	response := Post("/login/createSession", acc)

	var jwt login.JWT
	json.Unmarshal(response, &jwt)

	if jwt.Token != "" {
		t.Error("Expected session creation with incorrect username and password to return false, instead got true")
	}
}

func TestSessionStatus(t *testing.T) {
	response := Get("/sessions/validate")

	var result login.Status
	json.Unmarshal(response, &result)

	if result.Status != true {
		t.Error("Expected to validate session, instead failed to validate")
	}
}

func TestInvalidSession(t *testing.T) {
	response := Get("/sessions/validate", "abcdefghijklmnop")

	var result login.Status
	json.Unmarshal(response, &result)

	if result.Status != false {
		t.Error("Expected to fail to validate session, instead succeeded")
	}
}
