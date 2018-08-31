package xbl

import "net/url"

// ActivityStatus describes the current status of a Xbox LIVE gamertag.
type ActivityStatus struct {
	XUID    string           `json:"xuid"`
	State   string           `json:"state"`
	Devices []DeviceActivity `json:"devices,omitempty"`
}

type DeviceActivity struct {
	Type   string          `json:"type"`
	Titles []TitleActivity `json:"titles,omitempty"`
}

type TitleActivity struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Placement    string `json:"placement"`
	State        string `json:"state"`
	LastModified string `json:"lastModified"`
}

type activityStatusRequest struct {
	Users []string `json:"users"`
	Level string   `json:"level"`
}

// ActivityStatuses retrieves the activity statuses for the
// provided XIDs.
func (c *Client) ActivityStatuses(xids ...string) ([]ActivityStatus, error) {
	u := url.URL{
		Scheme: "https",
		Host:   "userpresence.xboxlive.com",
		Path:   "/users/batch",
	}

	var resp []ActivityStatus
	err := c.post(u.String(), vBoth, activityStatusRequest{
		Users: xids,
		Level: "all",
	}, &resp)
	return resp, err
}
