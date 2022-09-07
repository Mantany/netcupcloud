package session

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Define the base url & shop ID
const (
	shopNo  int    = 1
	baseUrl string = "https://www.netcup.eu/"
	timeout int    = 15
)

type EUShopSession struct {
	shopNo           int
	baseUrl          string
	customerNo       string
	customerPassword string
	isAuthenticated  bool
	httpClient       *http.Client
}

func NewEUShopSession(customer_no string, customer_password string) *EUShopSession {
	result := &EUShopSession{
		shopNo:           shopNo,
		baseUrl:          baseUrl,
		customerNo:       customer_no,
		customerPassword: customer_password,
		isAuthenticated:  false,
		httpClient:       &http.Client{},
	}
	return result
}

// Create a new Request to authenticate & set all nessesary header info for request
func (session *EUShopSession) newRequest(method string, path string, body io.Reader) (*http.Request, error) {
	url := baseUrl + path
	req, err := http.NewRequest(method, url, body)
	req.Header.Set("Host", "www.netcup.eu")
	req.Header.Set("Origin", "https://www.netcup.eu")
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8`)
	req.Header.Set("Accept-Language", `en-GB,en;q=0.5`)
	return req, err
}

func (session *EUShopSession) Do(req *http.Request) (*http.Response, error) {
	// TODO
	if !session.isAuthenticated {
		// authenticate
		session.auth()
	}
	resp, err := session.httpClient.Do(req)
	if err != nil {
		return nil, errors.New("Something went wrong while performing the request")
	}
	return resp, nil
}

func (session *EUShopSession) auth() error {
	// Build the authentication body - its a form, so:
	form := url.Values{}
	form.Add("kunden_laden", strconv.Itoa(session.shopNo))
	form.Add("knr", session.customerNo)
	form.Add("pwd", session.customerPassword)

	// Build the request and test it
	req, err := session.newRequest("POST", "bestellen/adresse.php", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		fmt.Println(err)
		return errors.New("Cant create the request.")
	}
	res, err := session.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return errors.New("Something went wrong while the Authentication request was submitted")
	}
	defer res.Body.Close()

	// Check the status code, must be 200!
	if res.StatusCode != 200 {
		fmt.Println("Faulty status code!", res.Status, res.StatusCode)
		return errors.New("Authentication Failed, bad Status Code")
	}

	// Parse the HTML Document, to find out if the login was successful:
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println(err)
		return errors.New("Something went wrong while trying parsing the authentication request.")
	}

	if doc.Find(".error").Length() != 0 {
		doc.Find(".error").Each(func(i int, s *goquery.Selection) {
			fmt.Println("ERROR During Authentication: ", s.Text())
		})
		return errors.New("EUShopSession - Error occured during Authentication")
	}

	if doc.Find("div:contains('Login successful')").Length() == 0 {
		fmt.Println("Authentication was not successful")
		return errors.New("EUShopSession - Authentication failed")
	}

	session.isAuthenticated = true
	return nil
}
