package xbl

import "testing"

func TestTitles(t *testing.T) {
	c := testClient(t)
	titles, err := c.Titles("2745051201447500")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", titles[0])
}
