package netcupcloud

import "github.com/mantany/netcupcloud/netcupcloud/session"

type Client struct {
	customerNo            string
	customerPassword      string
	scpUsername           string
	scpPassword           string
	scpWebservicePassword string
	euShopSession         session.EUShopSession
}

func NewClient(customerNo string, customerPassword string) *Client {

	client := &Client{
		customerNo:       customerNo,
		customerPassword: customerPassword,
	}

	client.euShopSession = *session.NewEUShopSession(client.customerNo, client.customerPassword)
	return client
}
