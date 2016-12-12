package xbl

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	credentials *credentials
}

// Expiry returns the time at which this client's credentials will
// no longer be valid and the client should be replaced with a new
// one obtained from Login.
func (c *Client) Expiry() time.Time {
	return c.credentials.expiresAt
}

type apiVersion int

const (
	vXbox360 apiVersion = 1
	vXboxOne apiVersion = 2
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
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", c.credentials.authHeader())
	req.Header.Set("x-xbl-contract-version", strconv.Itoa(int(v)))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(respBody)
}
