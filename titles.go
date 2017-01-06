package xbl

import (
	"fmt"
	"net/url"
	"time"
)

// Title represents an individual Xbox title.
type Title struct {
	ID                  int       `json:"titleId"`
	Type                int       `json:"titleType"`
	Platforms           []int     `json:"platforms"`
	Name                string    `json:"name"`
	LastPlayed          time.Time `json:"lastPlayed"`
	CurrentAchievements int       `json:"currentAchievements"`
	CurrentGamerscore   int       `json:"currentGamerscore"`
	Sequence            int       `json:"sequence"`
	TotalAchievements   int       `json:"totalAchievements"`
	TotalGamerscore     int       `json:"totalGamerscore"`
	RareUnlocks         []struct {
		RarityCategory   string `json:"rarityCategory"`
		NumUnlocks       int    `json:"numUnlocks"`
		IsRarestCategory bool   `json:"isRarestCategory"`
	} `json:"rareUnlocks"`
}

type titlesResponse struct {
	Titles     []Title    `json:"titles"`
	PagingInfo pagingInfo `json:"pagingInfo"`
}

// Titles retrieves all Xbox titles played by the provided XID.
func (c *Client) Titles(xid string, opts ...Option) ([]Title, error) {
	var reqOpts reqOptions
	for _, opt := range opts {
		opt(&reqOpts)
	}

	queryParams := url.Values{"orderBy": {"unlockTime"}}
	u := url.URL{
		Scheme:   "https",
		Host:     "achievements.xboxlive.com",
		Path:     fmt.Sprintf("/users/xuid(%s)/history/titles", xid),
		RawQuery: queryParams.Encode(),
	}

	var resp titlesResponse
	err := c.get(u.String(), vBoth, &resp)
	if err != nil {
		return nil, err
	}

	titles := make([]Title, 0, resp.PagingInfo.TotalRecords)
	for _, t := range titles {
		if t.LastPlayed.Before(reqOpts.updatedSince) {
			return titles, nil
		}
		titles = append(titles, t)
	}

	for resp.PagingInfo.ContinuationToken != nil {
		queryParams.Set("continuationToken", *resp.PagingInfo.ContinuationToken)
		u.RawQuery = queryParams.Encode()

		err := c.get(u.String(), vBoth, &resp)
		if err != nil {
			return nil, err
		}

		for _, t := range resp.Titles {
			if t.LastPlayed.Before(reqOpts.updatedSince) {
				return titles, nil
			}
			titles = append(titles, t)
		}
	}
	return titles, nil
}
