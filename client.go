package xbl

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

const userAgent = "SmartGlass/com.microsoft.smartglass (1610.1205.1554; OS Version 10.1.1 (Build 14B100))"

// Client encapsulates the entire Xbox Live API and a set
// of credentials to access the API.
//
// A Client is safe for concurrent access.
type Client struct {
	client      http.Client
	credentials *credentials
}

// Expiry returns the time at which this client's credentials will
// no longer be valid and the client should be replaced with a new
// one obtained from Login.
func (c *Client) Expiry() time.Time {
	return c.credentials.expiresAt
}

// UserID returns the XID user ID of the user who is authenticated with
// the API. All requests to the API are performed as this user.
func (c *Client) UserID() string {
	return c.credentials.xid
}

// Gamertag returns the gamertag of the user who is authenticated with the API.
func (c *Client) Gamertag() string {
	return c.credentials.gamertag
}

type apiVersion int

const (
	vXbox360 apiVersion = 1
	vXboxOne apiVersion = 2
	vBoth    apiVersion = 3
)

func (c *Client) post(url string, v apiVersion, body interface{}, respBody interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", c.credentials.authHeader())
	req.Header.Set("x-xbl-contract-version", strconv.Itoa(int(v)))

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(respBody)
}

func (c *Client) get(u string, v apiVersion, respBody interface{}) error {
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", c.credentials.authHeader())
	req.Header.Set("x-xbl-contract-version", strconv.Itoa(int(v)))

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(respBody)
}
