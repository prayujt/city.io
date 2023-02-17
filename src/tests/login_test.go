package tests

import (
	"api/login"

	"encoding/json"
	"fmt"
	"testing"
)

var session_id string

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

	var session login.Session
	json.Unmarshal(response, &session)

	if session.SessionId == "" {
		t.Error("Expected session creation to return true, instead got false")
	} else {
		session_id = session.SessionId
	}
}

func TestIncorrectUsername(t *testing.T) {
	acc := login.Account{
		Username: "test2",
		Password: "test",
	}

	response := Post("/login/createSession", acc)

	var session login.Session
	json.Unmarshal(response, &session)

	if session.SessionId != "" {
		t.Error("Expected session creation with incorrect username to return false, instead got true")
	}
}

func TestIncorrectPassword(t *testing.T) {
	acc := login.Account{
		Username: "test",
		Password: "test2",
	}

	response := Post("/login/createSession", acc)

	var session login.Session
	json.Unmarshal(response, &session)

	if session.SessionId != "" {
		t.Error("Expected session creation with incorrect password to return false, instead got true")
	}
}

func TestIncorrectUsernamePassword(t *testing.T) {
	acc := login.Account{
		Username: "testing",
		Password: "testing",
	}

	response := Post("/login/createSession", acc)

	var session login.Session
	json.Unmarshal(response, &session)

	if session.SessionId != "" {
		t.Error("Expected session creation with incorrect username and password to return false, instead got true")
	}
}

func TestSessionStatus(t *testing.T) {
	if session_id == "" {
		t.Error("Cannot test session status when session creation test failed")
		return
	}

	response := Get(fmt.Sprintf("/sessions/%s", session_id))

	var result login.Status
	json.Unmarshal(response, &result)

	if result.Status != true {
		t.Error("Expected to validate session, instead failed to validate")
	}
}

func TestInvalidSession(t *testing.T) {
	response := Get("/sessions/abcdefghijklmnop")

	var result login.Status
	json.Unmarshal(response, &result)

	if result.Status != false {
		t.Error("Expected to fail to validate session, instead succeeded")
	}
}

func TestSessionRemove(t *testing.T) {
	if session_id == "" {
		t.Error("Cannot test session status when session creation test failed")
		return
	}

	session := login.Session{
		SessionId: session_id,
	}

	response := Post("/sessions/logout", session)

	var result login.Status
	json.Unmarshal(response, &result)

	if result.Status != true {
		t.Error("Expected to succeed in logout, instead failed")
	}
}

func TestSessionStatusAfterLogout(t *testing.T) {
	if session_id == "" {
		t.Error("Cannot test session status when session creation test failed")
		return
	}

	response := Get(fmt.Sprintf("/sessions/%s", session_id))

	var result login.Status
	json.Unmarshal(response, &result)

	if result.Status != false {
		t.Error("Expected to fail to validate deleted session, instead succeeded")
	}
}
