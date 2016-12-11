package xbl

import "time"

type Client struct {
	credentials *credentials
}

// Expiry returns the time at which this client's credentials will
// no longer be valid and the client should be replaced with a new
// one obtained from Login.
func (c *Client) Expiry() time.Time {
	return c.credentials.expiresAt
}
