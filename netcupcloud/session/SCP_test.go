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
