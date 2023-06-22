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

type CCP struct {
	baseUrl         string
	ccpNo           string
	ccpPassword     string
	ccpMFASecret    string
	isAuthenticated bool
	httpClient      *http.Client
}

func NewCCP(ccpNo string, ccpPassword string, ccpMFASecret string) *CCP {
	jar, err := cookiejar.New(nil)
	if err != nil {
		fmt.Println("Error occured while creating the cookie jar")
	}
	client := &http.Client{Jar: jar}

	result := &CCP{
		baseUrl:         "https://www.customercontrolpanel.de/",
		ccpNo:           ccpNo,
		ccpPassword:     ccpPassword,
		ccpMFASecret:    ccpMFASecret,
		isAuthenticated: false,
		httpClient:      client,
	}
	return result
}

// Create a new Request to set all nessesary header info for request
func (session *CCP) newRequest(method string, path string, body io.Reader) (*http.Request, error) {
	url := session.baseUrl + path
	req, err := http.NewRequest(method, url, body)
	req.Header.Set("Host", "www.customercontrolpanel.de")
	req.Header.Set("Origin", "https://www.customercontrolpanel.de")
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8`)
	req.Header.Set("Accept-Language", `en-GB,en;q=0.5`)
	return req, err
}

func (session *CCP) auth() error {
	// request the site to get the cookie:
	pathCookie := "index.php?login_language=GB"
	// get the temp. site key & pickup right cookies:
	req, err := session.newRequest("GET", pathCookie, nil)
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

	//Request the form to authenticate:
	path := "start.php?login_language=GB"
	// Do the authentication request:
	form := url.Values{}
	form.Add("nocsrftoken", "")
	form.Add("action", "login")
	form.Add("sso", "")
	form.Add("ccp_user", session.ccpNo)
	form.Add("ccp_password", session.ccpPassword)

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
		fmt.Println("CCP login: wrong username or password")
		return errors.New("CCP authentication error: wrong username or password")
	}

	// check for multi-factor authentication:
	// currently not supported
	if authdoc.Find("input[name='tan']").Length() > 0 {
		fmt.Println("CCP login multi-factor authentication enabled")
		// if session.ccpMFASecret == "" {
		// 	return errors.New("CCP authentication error: multi-factor secret cant be null, please provide a secret")
		// }
		// totp := gotp.NewDefaultTOTP(session.ccpMFASecret)
		// totpForm := url.Values{}
		// totpForm.Add("nocsrftoken", "")
		// totpForm.Add("action", "login2fa")
		// totpForm.Add("sso", "")
		// totpForm.Add("tan", totp.Now())

		// authReq, err := session.newRequest("POST", path, strings.NewReader(form.Encode()))
		// if err != nil {
		// 	return err
		// }
		// authReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		// authRes, err := session.httpClient.Do(authReq)
		// if err != nil {
		// 	return err
		// }
		// if authRes.StatusCode != 200 {
		// 	fmt.Println("Faulty status code!", res.Status, res.StatusCode)
		// 	return errors.New("MFA Request failed, bad status code")
		// }

		// authdoc, err := goquery.NewDocumentFromReader(authRes.Body)
		// if err != nil {
		// 	return err
		// }
		fmt.Println("MFA not supported yet")
		return errors.New("MFA not supported yet")
	}

	// check if authentication was successful:
	testAuthReq, err := session.newRequest("GET", path, nil)
	if err != nil {
		return err
	}
	testAuthRes, err := session.httpClient.Do(testAuthReq)
	if err != nil {
		return err
	}
	if testAuthRes.StatusCode != 200 {
		fmt.Println("Faulty status code!", res.Status, res.StatusCode)
		return errors.New("request Failed, bad Status Code")
	}
	authdoc, err = goquery.NewDocumentFromReader(testAuthRes.Body)
	if err != nil {
		return err
	}

	if authdoc.Find("Angemeldet als").Length() > 0 || authdoc.Find("Logged in as").Length() > 0 {
		session.isAuthenticated = true
	} else {
		fmt.Println("CCP authentication went wrong!", authdoc.Text())
		return errors.New("CCP authentication failed")
	}
	return nil
}
