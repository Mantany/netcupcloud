package session

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/hooklift/gowsdl/soap"
)

// This test only succeed, if you have some active servers
// currently doesnt work, try using older version of https://github.com/hooklift/gowsdl/issues/232

func TestSCPSoap_GetVServers(t *testing.T) {
	headers := map[string]string{
		"User-Agent":   "Apache-HttpClient/4.5.5 (Java/16.0.1)",
		"Host":         "www.servercontrolpanel.de:443",
		"Content-Type": "text/xml;charset=UTF-8",
	}
	proxyUrl, err := url.Parse("http://127.0.0.1:8888")
	if err != nil {
		fmt.Println("lol")
	}
	myClient := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}

	soapClient := soap.NewClient("https://www.servercontrolpanel.de:443/SCP/WSEndUser", soap.WithHTTPHeaders(headers), soap.WithHTTPClient(myClient))

	scpSession := NewWSEndUser(soapClient)
	fmt.Println(scpSession)
	test, err := scpSession.GetVServers(&GetVServers{LoginName: "141116", Password: "VgZ3cw5PqV3F"})
	fmt.Println(err)
	fmt.Println(test)
	if err != nil {
		t.Error("Expected successful authentication")
	}

}
