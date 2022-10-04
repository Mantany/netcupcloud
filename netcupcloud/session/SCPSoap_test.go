package session

import (
	"fmt"
	"testing"
)

// This test only succeed, if you have some active servers
func TestSCPSoap_GetVServers(t *testing.T) {
	scpSession := NewSCPSoap(env.SCPNo, env.SCPSoapKey)
	l, err := scpSession.getVServers()
	if err != nil {
		t.Error("Expected successful authentication")
	}
	fmt.Print(l)
	t.Error("Expected successful authentication")

}
