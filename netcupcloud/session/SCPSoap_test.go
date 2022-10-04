package session

import (
	"fmt"
	"testing"
)

// This test only succeed, if you have some active servers
func TestSCP_ListAllServerst(t *testing.T) {
	scpSession := NewSCPSoap(env.SCPNo, env.SCPPwd)
	l, err := scpSession.listAllServer()
	if err != nil {
		t.Error("Expected successful authentication")
	}
	fmt.Print(l)
	t.Error("Expected successful authentication")

}
