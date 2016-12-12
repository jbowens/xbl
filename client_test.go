package xbl

import (
	"os"
	"testing"
)

func testClient(t *testing.T) *Client {
	user, pass := os.Getenv("XBL_USER"), os.Getenv("XBL_PASS")
	c, err := Login(user, pass)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestProfile(t *testing.T) {
	c := testClient(t)
	p, err := c.Profiles(c.credentials.xid)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", p)
}
