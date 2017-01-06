package xbl

import (
	"fmt"
	"net/url"
	"time"
)

// Achievement represents an individual user's achievement in a game.
type Achievement struct {
	ID                int       `json:"id"`
	TitleID           int       `json:"titleId"`
	Name              string    `json:"name"`
	Sequence          int       `json:"sequence"`
	Flags             int       `json:"flags"`
	UnlockedOnline    bool      `json:"unlockedOnline"`
	Unlocked          bool      `json:"unlocked"`
	IsSecret          bool      `json:"isSecret"`
	Platform          int       `json:"platform"`
	Gamescore         int       `json:"gamerscore"`
	Description       string    `json:"description"`
	LockedDescription string    `json:"lockedDescription"`
	Type              int       `json:"type"`
	TimeUnlocked      time.Time `json:"timeUnlocked"`
	Rarity            Rarity    `json:"rarity"`
}

// Rarity describes the rarity of a particular achievement.
type Rarity struct {
	CurrentCategory   string `json:"currentCategory"`
	CurrentPercentage int    `json:"currentPercentage"`
}

type achievementsResponse struct {
	Achievements []Achievement `json:"achievements"`
	PagingInfo   pagingInfo    `json:"pagingInfo"`
}

type pagingInfo struct {
	ContinuationToken *string `json:"continuationToken"`
	TotalRecords      uint64  `json:"totalRecords"`
}

// Achievements retrieves all achievements for the provided XID.
func (c *Client) Achievements(xid string) ([]Achievement, error) {
	queryParams := url.Values{"maxItems": {"1000"}, "orderBy": {"EndingSoon"}}
	u := url.URL{
		Scheme:   "https",
		Host:     "achievements.xboxlive.com",
		Path:     fmt.Sprintf("/users/xuid(%s)/achievements", xid),
		RawQuery: queryParams.Encode(),
	}

	var resp achievementsResponse
	err := c.get(u.String(), vBoth, &resp)
	if err != nil {
		return nil, err
	}

	achievements := make([]Achievement, 0, resp.PagingInfo.TotalRecords)
	achievements = append(achievements, resp.Achievements...)
	for resp.PagingInfo.ContinuationToken != nil {
		queryParams.Set("continuationToken", *resp.PagingInfo.ContinuationToken)
		u.RawQuery = queryParams.Encode()

		err := c.get(u.String(), vBoth, &resp)
		if err != nil {
			return nil, err
		}
		achievements = append(achievements, resp.Achievements...)
	}
	return achievements, nil
}
