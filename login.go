package xbl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// Login logs into the Microsoft Live API, returning a client with a
// fresh access token.
func Login(username, password string) (*Client, error) {
	c, err := loginHTTPClient()
	if err != nil {
		return nil, err
	}

	authURL, token, err := oauthAuthorizeURL(c)
	if err != nil {
		return nil, err
	}
	creds, err := getCredentials(c, username, password, authURL, token)
	if err != nil {
		return nil, err
	}

	err = authenticate(c, creds)
	if err != nil {
		return nil, err
	}
	err = authorize(c, creds)
	if err != nil {
		return nil, err
	}

	return &Client{
		credentials: creds,
		client:      *http.DefaultClient,
	}, nil
}

type credentials struct {
	token        string
	uhs          string
	xid          string
	gamertag     string
	accessToken  string
	refreshToken string
	userID       string
	expiresAt    time.Time
}

func (c *credentials) authHeader() string {
	return fmt.Sprintf("XBL3.0 x=%s;%s", c.uhs, c.token)
}

func loginHTTPClient() (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	return &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // never follow redirects
		},
		Jar: jar,
	}, nil
}

var sFFTagRegexp = regexp.MustCompile("sFTTag:.*value=\"(.*)\"/>")
var authorizeURLRegexp = regexp.MustCompile("urlPost:'([^']+)'")

func oauthAuthorizeURL(c *http.Client) (string, string, error) {
	var defaultOauthAuthorizeRequest = url.Values{
		"client_id":     {"0000000048093EE3"},
		"redirect_uri":  {"https://login.live.com/oauth20_desktop.srf"},
		"response_type": {"token"},
		"display":       {"touch"},
		"scope":         {"service::user.auth.xboxlive.com::MBI_SSL"},
		"locale":        {"en"},
	}
	u := url.URL{
		Scheme:   "https",
		Host:     "login.live.com",
		Path:     "/oauth20_authorize.srf",
		RawQuery: defaultOauthAuthorizeRequest.Encode(),
	}
	resp, err := c.Get(u.String())
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}
	matches := sFFTagRegexp.FindStringSubmatch(string(body))
	if len(matches) < 2 {
		return "", "", errors.New("unable to find sFFTag")
	}
	sFFTag := matches[1] // 1st submatch
	matches = authorizeURLRegexp.FindStringSubmatch(string(body))
	if len(matches) < 2 {
		return "", "", errors.New("unable to find authorize URL")
	}
	authorizeURL := matches[1]

	return authorizeURL, sFFTag, nil
}

func getCredentials(c *http.Client, username, password, authorizeURL, ppft string) (*credentials, error) {
	p := url.Values{
		"login":        {username},
		"loginfmt":     {username},
		"passwd":       {password},
		"PPFT":         {ppft},
		"PPSX":         {"Passp"},
		"ps":           {"2"},
		"canary":       {},
		"ctx":          {},
		"SI":           {"Sign in"},
		"type":         {"11"},
		"fspost":       {"0"},
		"NewUser":      {"1"},
		"LoginOptions": {"1"},
		"i2":           {"39"},
		"i3":           {"36728"},
		"m1":           {"768"},
		"m2":           {"1184"},
		"m3":           {"0"},
		"i12":          {"1"},
		"i17":          {"0"},
		"i18":          {"__DefaultLoginPaginatedStrings|1,__DefaultLogin_PCore|1"},
		"i21":          {"0"},
	}

	req, err := http.NewRequest("POST", authorizeURL, strings.NewReader(p.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", authorizeURL)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.98 Safari/537.36")
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 302 {
		return nil, errors.New("unable to login; verify your username and password")
	}
	u, err := url.Parse(resp.Header.Get("Location"))
	if err != nil {
		return nil, err
	}

	authValues, err := url.ParseQuery(u.Fragment)
	if err != nil {
		return nil, err
	}
	if len(authValues["access_token"]) < 1 {
		return nil, fmt.Errorf("bad login access token, url: %s", u)
	}

	return &credentials{
		accessToken:  authValues["access_token"][0],
		refreshToken: authValues["refresh_token"][0],
		userID:       authValues["user_id"][0],
	}, nil
}

func authenticate(c *http.Client, creds *credentials) error {
	const authenticateURL = "https://user.auth.xboxlive.com/user/authenticate"
	body, err := json.Marshal(map[string]interface{}{
		"RelyingParty": "http://auth.xboxlive.com",
		"TokenType":    "JWT",
		"Properties": map[string]interface{}{
			"AuthMethod": "RPS",
			"SiteName":   "user.auth.xboxlive.com",
			"RpsTicket":  creds.accessToken,
		},
	})
	if err != nil {
		return err
	}

	resp, err := c.Post(authenticateURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var authResponse struct {
		Token         string `json:"Token"`
		DisplayClaims struct {
			XUI []struct {
				UHS string `json:"uhs"`
			} `json:"xui"`
		} `json:"DisplayClaims"`
	}
	err = json.NewDecoder(resp.Body).Decode(&authResponse)
	if err != nil {
		return err
	}

	creds.token = authResponse.Token
	creds.uhs = authResponse.DisplayClaims.XUI[0].UHS
	return nil
}

func authorize(c *http.Client, creds *credentials) error {
	const authorizeURL = "https://xsts.auth.xboxlive.com/xsts/authorize"
	body, err := json.Marshal(map[string]interface{}{
		"RelyingParty": "http://xboxlive.com",
		"TokenType":    "JWT",
		"Properties": map[string]interface{}{
			"UserTokens": []string{creds.token},
			"SandboxId":  "RETAIL",
		},
	})
	if err != nil {
		return err
	}

	resp, err := c.Post(authorizeURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var authorizeResponse struct {
		Token         string    `json:"Token"`
		NotAfter      time.Time `json:"NotAfter"`
		DisplayClaims struct {
			XUI []struct {
				XID      string `json:"xid"`
				Gamertag string `json:"gtg"`
			} `json:"xui"`
		} `json:"DisplayClaims"`
	}
	err = json.NewDecoder(resp.Body).Decode(&authorizeResponse)
	if err != nil {
		return err
	}

	creds.token = authorizeResponse.Token
	creds.expiresAt = authorizeResponse.NotAfter
	creds.xid = authorizeResponse.DisplayClaims.XUI[0].XID
	creds.gamertag = authorizeResponse.DisplayClaims.XUI[0].Gamertag
	return nil
}
