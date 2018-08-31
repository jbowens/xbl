package xbl

import "testing"

func TestActivity(t *testing.T) {
	c := testClient(t)
	activities, err := c.ActivityStatuses("2745051201447500")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", activities[0])
}
