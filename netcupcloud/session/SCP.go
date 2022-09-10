package session

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
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
		return nil, errors.New("request Failed, bad Status Code")
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
	// TODO check the website validation

	// Set the cookie header
	//urltest, err := url.Parse("https://www.servercontrolpanel.de/")
	//req.Header.Set("Cookie", session.httpClient.Jar.Cookies(urltest)[0].String())

	// renew the site key, because it changes after the request
	session.renewSiteKey(*authdoc)

	if session.scpSiteKey != "" {
		session.isAuthenticated = true
	}
	return nil
}

func (session *SCP) renewSiteKey(doc goquery.Document) error {
	keys, err := getJavaScriptVarsFromHTML(doc, "site_key")
	if err != nil {
		return err
	}
	if len(keys) == 0 {
		fmt.Println("cant retrieve sitekey")
		return errors.New("cant retrieve sitekey")
	}
	session.scpSiteKey = keys[0]["site_key"]
	return nil
}

// TODO Should probably move this one to a seperate package
// This method is used to get all java script vars & arrays out of an html document
// @param: html: goquery Document
// @param: var_name: the name of the java-script var
// 	i = {} -> Match!
// 	i['test'] = x

// [map[site_key:xck5GEb6t5y4ehd2tp3zEBpWxgn5H7Wg]]
// [{i: {}, i['test']: x}]
func getJavaScriptVarsFromHTML(doc goquery.Document, varName string) ([]map[string]string, error) {
	//Grab the whole java script code from the site:
	t := strings.Join(doc.Find("script").Map(func(i int, s *goquery.Selection) string { return s.Text() }), " ")
	var result = []map[string]string{}

	// Find all java script variables:
	pattern := regexp.MustCompile(`(?P<key>\S*?)[^\S\r\n]*?=[^\S\r\n]*?[\'\"]?(?P<value>\S*?)[\'\"]?;`)
	search := pattern.FindAllStringSubmatch(t, -1)

	// this pattern is used to find arrays
	// f.e: links['hello'] = x
	patternFindArrayMatch, err := regexp.Compile(varName + `\[[\'"]?(\S*?)["\']?\]`)
	if err != nil {
		return nil, err
	}
	// now find the right name, and safe the matches!
	for _, s := range search {
		if len(s) == 3 {
			// search for array match:
			search := patternFindArrayMatch.FindString(s[1])
			if s[1] == varName || search != "" {
				result = append(result, map[string]string{s[1]: s[2]})
			}
		}
	}
	return result, nil
}

// uses the SCP to list all server, returns the scpID and rdns of the server
func (session *SCP) ListAllServersWithID() ([]map[string]string, error) {
	if session.scpSiteKey == "" || !session.isAuthenticated {
		session.auth()
	}
	path := "SCP/Home"
	form := url.Values{}
	form.Add("site_key", session.scpSiteKey)
	form.Add("statusboxText", "no+new+status+available")

	req, err := session.newRequest("POST", path, strings.NewReader(form.Encode()))
	req.Header.Set("Referer", "https://www.servercontrolpanel.de/SCP/Home")
	//important header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		return nil, err
	}
	res, err := session.do(req)

	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}
	links, err := getJavaScriptVarsFromHTML(*doc, "links")
	if err != nil {
		return nil, err
	}
	if len(links) == 0 {
		return nil, nil
	}

	var result = []map[string]string{}

	for _, link := range links {
		for key, value := range link {
			rdnsPattern := regexp.MustCompile(`links\[\'(?P<link>.*)\'\]`)
			rdns := rdnsPattern.FindStringSubmatch(key)
			idPattern := regexp.MustCompile(`VServersKVM\?selectedVServerId=([0-9]*)`)
			id := idPattern.FindStringSubmatch(value)
			if len(rdns) == 2 && len(id) == 2 && rdns[1] != "" && id[1] != "" {
				result = append(result, map[string]string{"rDNS": rdns[1], "scpID": id[1]})
			}
		}
	}
	return result, nil
}
