package session

import (
	"testing"
)

func TestCCP_AuthWrongPasswordWrongUsernameWrongMFA(t *testing.T) {
	ccpSession := NewCCP("Wrong", "password", "mfaSecret")
	err := ccpSession.auth()
	if err == nil {
		t.Error("Expected wrong password error!")
	}
}

func TestCCP_AuthRightPasswordRightUsernameRightMFA(t *testing.T) {
	ccpSession := NewCCP(env.CCPNo, env.CCPPwd, env.CCPMFASecret)
	err := ccpSession.auth()
	if err != nil {
		t.Error("Expected successful authentication, please make sure that you have set the environment variables correctly")
	}
}

// For this test to work, you need to have have MFA enabled for your account
func TestCCP_AuthRightPasswordRightUsernameWrongMFA(t *testing.T) {
	ccpSession := NewCCP(env.CCPNo, env.CCPPwd, "NEWCOOL2FAKEYAAW")
	err := ccpSession.auth()
	if err == nil {
		t.Error("Expected wrong MFA error!")
	}
}
