package session

import (
	"context"
	"fmt"
	"net/http"

	"github.com/globusdigital/soap"
)

type SCPSoap struct {
	wsdlUrl         string
	scpSoapUsername string
	scpSoapKey      string
	soapClient      *soap.Client
}

func NewSCPSoap(scpSoapUsername string, scpSoapKey string) *SCPSoap {

	soapClient := soap.NewClient("https://www.servercontrolpanel.de/WSEndUser?wsdl", &soap.BasicAuth{scpSoapUsername, scpSoapKey})
	soapClient.UseSoap12()
	soapClient.ContentType = "application/xml"
	soapClient.UserAgent = `Mozilla/5.0`
	result := &SCPSoap{
		wsdlUrl:         "https://www.servercontrolpanel.de/WSEndUser?wsdl",
		scpSoapUsername: scpSoapUsername,
		scpSoapKey:      scpSoapKey,
		soapClient:      soapClient,
	}
	return result
}

// FooRequest a simple request
type getVServersRequest struct {
}

// FooResponse a simple response
type getVServersResponse struct {
	servers []string
}

func (session *SCPSoap) getVServers() (http.Response, error) {
	response := &getVServersResponse{}
	httpResponse, err := session.soapClient.Call(context.TODO(), "getVServers", &getVServersRequest{}, response)

	fmt.Println("Hello")
	fmt.Println(err)

	if err != nil {
		return http.Response{}, err
	}
	fmt.Println(httpResponse)
	return http.Response{}, nil
}
