package session

import (
	"context"
	"encoding/xml"
	"log"
	"net/http"

	"github.com/globusdigital/soap"
)

type SCPSoap struct {
	wsdlUrl         string
	scpSoapUsername string
	scpSoapPassword string
	soapClient      *soap.Client
}

func NewSCPSoap(scpSoapUsername string, scpSoapPassword string) *SCPSoap {

	soapClient := soap.NewClient("https://www.servercontrolpanel.de/WSEndUser?wsdl", &soap.BasicAuth{scpSoapUsername, scpSoapPassword})
	result := &SCPSoap{
		wsdlUrl:         "https://www.servercontrolpanel.de/WSEndUser?wsdl",
		scpSoapUsername: scpSoapUsername,
		scpSoapPassword: scpSoapPassword,
		soapClient:      soapClient,
	}
	return result
}

// FooRequest a simple request
type FooRequest struct {
	XMLName xml.Name `xml:"fooRequest"`
	Foo     string
}

// FooResponse a simple response
type FooResponse struct {
	Bar string
}

func (session *SCPSoap) listAllServer() (http.Response, error) {
	response := &FooResponse{}
	httpResponse, err := session.soapClient.Call(context.Context, "getVServers", &FooRequest{Foo: "hello i am foo"}, response)
	if err != nil {
		return http.Response{}, err
	}

	log.Println(response.Bar, httpResponse.Status)
	return http.Response{}, nil
}
