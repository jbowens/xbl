package xbl

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// Profile describe's a gamer's Xbox Live account.
type Profile struct {
	ID             string `json:"id"`
	Gamertag       string `json:"gamertag,omitempty"`
	Gamerscore     int    `json:"gamerscore,omitempty"`
	GamerPicture   string `json:"picture_url,omitempty"`
	AccountTier    string `json:"account_tier,omitempty"`
	Reputation     string `json:"reputation,omitempty"`
	PreferredColor string `json:"-"`
	Tenure         int    `json:"-"`
	Bio            string `json:"bio,omitempty"`
	Location       string `json:"location,omitempty"`
}

var profileSettings = []string{
	"Gamerscore", "Gamertag", "GameDisplayPicRaw", "AccountTier", "Bio",
	"XboxOneRep", "PreferredColor", "TenureLevel", "Location",
}

// Profile retrieves the profile for the provided gamertag.
func (c *Client) Profile(gamertag string) (*Profile, error) {
	u := url.URL{
		Scheme:   "https",
		Host:     "profile.xboxlive.com",
		Path:     fmt.Sprintf("/users/gt(%s)/settings", gamertag),
		RawQuery: url.Values{"settings": {strings.Join(profileSettings, ",")}}.Encode(),
	}

	var resp profileResponse
	err := c.get(u.String(), vBoth, &resp)
	if err != nil {
		return nil, err
	}

	profiles, err := resp.parseProfiles()
	if err != nil {
		return nil, err
	}
	if len(profiles) < 1 {
		return nil, nil
	}
	return profiles[0], nil
}

// Profiles retrieves the profiles for the users with the provided IDs.
func (c *Client) Profiles(userIDs ...string) ([]*Profile, error) {
	const url = "https://profile.xboxlive.com/users/batch/profile/settings"

	req := struct {
		Settings []string `json:"settings"`
		UserIDs  []string `json:"userIds"`
	}{
		Settings: profileSettings,
		UserIDs:  userIDs,
	}
	var resp profileResponse
	err := c.post(url, vBoth, req, &resp)
	if err != nil {
		return nil, err
	}

	profiles, err := resp.parseProfiles()
	if err != nil {
		return nil, err
	}
	return profiles, nil
}

type profileResponse struct {
	ProfileUsers []struct {
		ID       string              `json:"id"`
		Settings []map[string]string `json:"settings"`
	} `json:"profileUsers"`
}

func (resp *profileResponse) parseProfiles() ([]*Profile, error) {
	profiles := make([]*Profile, 0, len(resp.ProfileUsers))
	for _, pu := range resp.ProfileUsers {
		m := map[string]string{}
		for _, setting := range pu.Settings {
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

		profiles = append(profiles, &Profile{
			ID:             pu.ID,
			Gamertag:       m["Gamertag"],
			Gamerscore:     gamerscore,
			GamerPicture:   m["GameDisplayPicRaw"],
			AccountTier:    m["AccountTier"],
			Reputation:     m["XboxOneRep"],
			PreferredColor: m["PreferredColor"],
			Tenure:         tenure,
			Bio:            m["Bio"],
			Location:       m["Location"],
		})
	}
	return profiles, nil
}
