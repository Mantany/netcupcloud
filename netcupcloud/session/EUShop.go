package session

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	timeout int = 15
)

type EUShop struct {
	shopNo           int
	baseUrl          string
	customerNo       string
	customerPassword string
	isAuthenticated  bool
	httpClient       *http.Client
}

// Define the base url & shop ID
func NewEUShop(customerNo string, customerPassword string) *EUShop {
	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Println("Error creating cookie jar")
	}
	client := &http.Client{
		Jar: jar,
	}

	result := &EUShop{
		shopNo:           1,
		baseUrl:          "https://www.netcup.eu/",
		customerNo:       customerNo,
		customerPassword: customerPassword,
		isAuthenticated:  false,
		httpClient:       client,
	}
	return result
}

// Create a new Request to authenticate & set all nessesary header info for request
func (session *EUShop) newRequest(method string, path string, body io.Reader) (*http.Request, error) {
	url := session.baseUrl + path
	req, err := http.NewRequest(method, url, body)
	req.Header.Set("Host", "www.netcup.eu")
	req.Header.Set("Origin", "https://www.netcup.eu")
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8`)
	req.Header.Set("Accept-Language", `en-GB,en;q=0.5`)
	return req, err
}

func (session *EUShop) do(req *http.Request) (*http.Response, error) {

	if !session.isAuthenticated {
		// authenticate
		session.auth()
	}
	resp, err := session.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		fmt.Println("Faulty status code!", resp.Status, resp.StatusCode)
		return nil, errors.New("authentication Failed, bad Status Code")
	}
	return resp, nil
}

func (session *EUShop) auth() error {
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
		return errors.New("cant create the request")
	}
	res, err := session.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return errors.New("something went wrong while the Authentication request was submitted")
	}
	defer res.Body.Close()

	// Check the status code, must be 200!
	if res.StatusCode != 200 {
		fmt.Println("Faulty status code!", res.Status, res.StatusCode)
		return errors.New("authentication Failed, bad Status Code")
	}

	// Parse the HTML Document, to find out if the login was successful:
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
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

// puts a product with specified ID into the chart
func (session *EUShop) PutIntoChart(id int) error {
	path := "bestellen/warenkorb_add.php?produkt=" + strconv.Itoa(id)
	req, err := session.newRequest("GET", path, nil)
	if err != nil {
		fmt.Println("PutIntoChart", err)
		return err
	}
	res, err := session.do(req)
	if err != nil {
		fmt.Println("PutIntoChart", err)
		return err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("PutIntoChart", err)
		return err
	}
	if doc.Find("main:contains('Error')").Length() != 0 {
		doc.Find("main:contains('Error')").Each(func(i int, s *goquery.Selection) {
			fmt.Println("error during PutIntoChart", s.Text())
		})
		return errors.New("EUShopSession - Error occured during PutIntoChart")
	}
	if doc.Find("h1:contains('The product was added to cart.')").Length() == 0 {
		fmt.Println("Put Into Chart was not successful  - Website didnt show confirmation")
		return errors.New("PutIntoChart - Website didnt show confirmation")
	}
	return nil
}

// Releases the order -> will throw an error if you dont have something in your chart
func (session *EUShop) ReleaseOrder() error {
	path := "bestellen/bestellung_ausfuehren.php"
	form := url.Values{}
	form.Add("agb", "1")
	form.Add("widerruf_gelesen", "1")

	req, err := session.newRequest("POST", path, strings.NewReader(form.Encode()))
	if err != nil {
		fmt.Println("releaseOrder")
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "https://www.netcup.eu/bestellen/bestellbestaetigung.php")

	res, err := session.do(req)
	if err != nil {
		fmt.Println("releaseOrder")
		return err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("releaseOrder")
		return err
	}

	if doc.Find("span:contains('To make an order, there must be at least one product in your cart.')").Length() != 0 {
		fmt.Println("ReleaseOrder was not successful - no element in chart found")
		return errors.New("ReleaseOrder - No element in chart found")
	}

	if doc.Find("p:contains('Thank you for your purchase at netcup!')").Length() == 0 {
		fmt.Println("ReleaseOrder no Order confirmation given!")
		return errors.New("ReleaseOrder - no Order confirmation given")
	}
	return nil
}
