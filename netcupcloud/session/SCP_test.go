package session

import (
	"testing"
)

func TestSCP_AuthWrongPasswordWrongUsername(t *testing.T) {
	scpSession := NewSCP("Wrong", "password")
	err := scpSession.auth()
	if err == nil {
		t.Error("Expected wrong password error!")
	}
}

func TestSCP_AuthRightPasswordRightUsername(t *testing.T) {
	scpSession := NewSCP(env.SCPNo, env.SCPPwd)
	err := scpSession.auth()
	if err != nil {
		t.Error("Expected successful authentication")
	}
}

// This test only succeed, if you have some active servers
func TestSCP_ListAllServersByIdWithServerAvailable(t *testing.T) {
	scpSession := NewSCP(env.SCPNo, env.SCPPwd)
	l, err := scpSession.ListAllServersWithID()
	if l == nil || err != nil || len(l) == 0 {
		t.Error("Expected successful listing of available server")
	}
}
