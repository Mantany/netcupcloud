package session

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type SCP struct {
	baseUrl         string
	scpNo           string
	scpPassword     string
	scpSiteKey      string
	isAuthenticated bool
	httpClient      *http.Client
}

func NewSCP(scpNo string, scpPassword string) *SCP {
	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Println("Error occured while creating the cookie jar")
	}
	client := &http.Client{Jar: jar}

	result := &SCP{
		baseUrl:     "https://www.servercontrolpanel.de/",
		scpNo:       scpNo,
		scpPassword: scpPassword,
		httpClient:  client,
	}
	return result
}

// Create a new Request to authenticate & set all nessesary header info for request
func (session *SCP) newRequest(method string, path string, body io.Reader) (*http.Request, error) {
	url := session.baseUrl + path
	req, err := http.NewRequest(method, url, body)
	req.Header.Set("Host", "www.servercontrolpanel.de")
	req.Header.Set("Origin", "https://www.servercontrolpanel.de")
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8`)
	req.Header.Set("Accept-Language", `en-GB,en;q=0.5`)
	return req, err
}

func (session *SCP) do(req *http.Request) (*http.Response, error) {
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
		return nil, errors.New("Request Failed, bad Status Code")
	}
	return resp, nil
}

func (session *SCP) auth() error {
	path := "SCP/Login"
	// get the temp. site key & pickup right cookies:
	req, err := session.newRequest("GET", path, nil)
	if err != nil {
		return err
	}
	fmt.Println(req.Body, req.Header)
	res, err := session.httpClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		fmt.Println("Faulty status code!", res.Status, res.StatusCode)
		return errors.New("request Failed, bad Status Code")
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}
	siteKey, t := doc.Find("#site_key").Attr("value")
	if !t {
		fmt.Println("Couldnt retrieve the site key")
		return errors.New("authentication Failed, cant retrieve the site key")
	}

	// Do the authentication request:
	form := url.Values{}
	form.Add("site_key", siteKey)
	form.Add("username", session.scpNo)
	form.Add("password", session.scpPassword)

	authReq, err := session.newRequest("POST", path, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	authReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	authRes, err := session.httpClient.Do(authReq)
	if err != nil {
		return err
	}
	if authRes.StatusCode != 200 {
		fmt.Println("Faulty status code!", res.Status, res.StatusCode)
		return errors.New("request Failed, bad Status Code")
	}
	authdoc, err := goquery.NewDocumentFromReader(authRes.Body)
	if err != nil {
		return err
	}

	if authdoc.Find(".alert:contains('Login failed! Please check your username and password.')").Length() > 0 {
		fmt.Println("SCP login: wrong username or password")
		return errors.New("SCP authentication error: wrong username or password")
	}

	if authdoc.Find(".alert").Length() > 0 {
		authdoc.Find(".alert").Each(func(i int, s *goquery.Selection) {
			fmt.Println("SCP authentication went wrong!", s.Text())
		})
		return errors.New("SCP authentication failed")
	}

	fmt.Println(authdoc.Html())

	return nil
}

func (session *SCP) renewSiteKey(htmlBody io.Reader) error {
	keys := getJavaScriptVarsFromHTML(htmlBody, "site_key", false)
	if len(keys) == 0 {
		fmt.Println("Cant retrieve sitekey!")
		return errors.New("Cant retrieve sitekey!")
	}
	session.scpSiteKey = keys[0]["site_key"]
	return nil
}

// This method is used to get all java script vars out of an html document
// @param: html: request body
// @param: var_name: the name of the java-script var
// @param: match_arrays: wether arrays should be matched f.e.
// 	i = {} -> Match!
// 	i['test'] = x

// Output with match_array = True
// [{i: {}, i['test']: x}]
func getJavaScriptVarsFromHTML(htmlBody io.Reader, var_name string, match_arrays bool) []map[string]string {

}

func (session *SCP) ListAllServerWithID() {

}
