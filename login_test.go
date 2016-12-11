package xbl

import "testing"

func TestOAuthAuthorizeURL(t *testing.T) {
	c, err := loginHTTPClient()
	if err != nil {
		t.Fatal(err)
	}

	u, tag, err := oauthAuthorizeURL(c)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(u)
	t.Log(tag)
}
