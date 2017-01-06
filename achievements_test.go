package xbl

import "testing"

func TestAchievements(t *testing.T) {
	c := testClient(t)
	achievements, err := c.Achievements("2533274792255104")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%#v", achievements[0])
}
