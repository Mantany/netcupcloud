package session

import (
	"github.com/hooklift/gowsdl/soap"
)

type SCPSoap struct {
	wsdlUrl         string
	scpSoapUsername string
	scpSoapKey      string
	soapClient      *soap.Client
}

// FooRequest a simple request
type getVServersRequest struct {
	loginName string
	password  string
}

// will hold the complete xml response:
type getVServersResponse struct {
	getVServersResponse string `xml:"getVServersResponse"`
}

// will hold the right data struct, of the xml response object:
type getVServersResult struct {
	response string `xml:"return"`
}

var (
	r getVServersResponse
)

type getVServerStateRequest struct {
	vserverName string
}

type getVServersStateResponse struct {
}
