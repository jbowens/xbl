package xbl

import (
	"fmt"
	"strconv"
)

// Profile describe's a gamer's Xbox Live account.
type Profile struct {
	ID             string `json:"-"`
	Gamertag       string `json:"Gamertag"`
	Gamerscore     int    `json:"Gamerscore"`
	GamerPicture   string `json:"GameDisplayPicRaw"`
	AccountTier    string `json:"AccountTier"`
	Reputation     string `json:"XboxOneRep"`
	PreferredColor string `json:"PreferredColor"`
	Tenure         int    `json:"TenureLevel"`
}

// Profile retrieves the profile for a user with the provided ID.
func (c *Client) Profile(userID string) (*Profile, error) {
	const url = "https://profile.xboxlive.com/users/batch/profile/settings"

	req := struct {
		Settings []string `json:"settings"`
		UserIDs  []string `json:"userIds"`
	}{
		Settings: []string{
			"Gamerscore", "Gamertag", "GameDisplayPicRaw", "AccountTier",
			"XboxOneRep", "PreferredColor", "TenureLevel",
		},
		UserIDs: []string{userID},
	}

	var resp struct {
		ProfileUsers []struct {
			ID       string              `json:"id"`
			Settings []map[string]string `json:"settings"`
		} `json:"profileUsers"`
	}
	err := c.post(url, vXboxOne, req, &resp)
	if err != nil {
		return nil, err
	}
	if len(resp.ProfileUsers) == 0 {
		return nil, fmt.Errorf("user id %s not found", userID)
	}
	m := map[string]string{}
	for _, setting := range resp.ProfileUsers[0].Settings {
		m[setting["id"]] = setting["value"]
	}
	tenure, err := strconv.Atoi(m["TenureLevel"])
	if err != nil {
		return nil, err
	}
	gamerscore, err := strconv.Atoi(m["Gamerscore"])
	if err != nil {
		return nil, err
	}

	return &Profile{
		ID:             userID,
		Gamertag:       m["Gamertag"],
		Gamerscore:     gamerscore,
		GamerPicture:   m["GameDisplayPicRaw"],
		AccountTier:    m["AccountTier"],
		Reputation:     m["XboxOneRep"],
		PreferredColor: m["PreferredColor"],
		Tenure:         tenure,
	}, nil
}
