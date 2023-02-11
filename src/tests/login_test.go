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

func TestAccountLogin(t *testing.T) {
	acc := login.Account{
		Username: "test",
		Password: "test",
	}

	body := Post("/login/createSession", acc)

	var session login.Session
	json.Unmarshal(body, &session)

	if !session.Status {
		t.Error("Expected session to return true, instead got false")
	}
}
