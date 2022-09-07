package netcupcloud

import "github.com/Mantany/netcupcloud/netcupcloud/session"

type Client struct {
	customer_no             string
	customer_password       string
	scp_username            string
	scp_password            string
	scp_webservice_password string
	eu_shop_session         session.EUShopSession
}

func NewClient(customer_no string, customer_password string) *Client {

	client := &Client{
		customer_no:       customer_no,
		customer_password: customer_password,
	}

	client.eu_shop_session = *session.NewEUShopSession(client.customer_no, client.customer_password)

	return client
}
