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

	body := Post("/login/createAccount", acc)

	if string(body) != "true" {
		t.Error("Expected account creation to return true, instead got false")
	}
}

func TestDuplicateAccount(t *testing.T) {
	acc := login.Account{
		Username: "test",
		Password: "test",
	}

	body := Post("/login/createAccount", acc)

	if string(body) != "false" {
		t.Error("Expected duplicate account creation to return false, instead got true")
	}
}

func TestAccountLogin(t *testing.T) {
	acc := login.Account{
		Username: "test",
		Password: "test",
	}

	body := Post("/login/createSession", acc)

	var session login.Session
	json.Unmarshal(body, &session)

	if !session.Status {
		t.Error("Expected session creation to return true, instead got false")
	}
}

func TestIncorrectUsername(t *testing.T) {
	acc := login.Account{
		Username: "test2",
		Password: "test",
	}

	body := Post("/login/createSession", acc)

	var session login.Session
	json.Unmarshal(body, &session)

	if session.Status {
		t.Error("Expected session creation with incorrect username to return false, instead got true")
	}
}

func TestIncorrectPassword(t *testing.T) {
	acc := login.Account{
		Username: "test",
		Password: "test2",
	}

	body := Post("/login/createSession", acc)

	var session login.Session
	json.Unmarshal(body, &session)

	if session.Status {
		t.Error("Expected session creation with incorrect password to return false, instead got true")
	}
}

func TestIncorrectUsernamePassword(t *testing.T) {
	acc := login.Account{
		Username: "testing",
		Password: "testing",
	}

	body := Post("/login/createSession", acc)

	var session login.Session
	json.Unmarshal(body, &session)

	if session.Status {
		t.Error("Expected session creation with incorrect username and password to return false, instead got true")
	}
}
