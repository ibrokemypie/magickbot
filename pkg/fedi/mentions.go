package fedi

import (
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

// Notification - Mastodon notification object
type Notification struct {
	ID      string  `json:"id"`
	Type    string  `json:"type"`
	Account Account `json:"account"`
	Status  Status  `json:"status"`
}

// GetMentions - Get the bots mention notifications
func GetMentions(instanceURL, accessToken string) ([]Notification, error) {
	mentions := make([]Notification, 0)
	lastID := viper.GetString("last_mention_id")

	u, err := url.Parse(instanceURL + "/api/v1/notifications")
	if err != nil {
		return nil, err
	}

	_, err = resty.New().R().
		SetAuthToken(accessToken).
		SetQueryParamsFromValues(url.Values{
			"min_id":          []string{lastID},
			"exclude_types[]": []string{"follow", "favourite", "reblog", "poll", "follow_request"},
		}).
		SetResult(&mentions).
		Get(u.String())
	if err != nil {
		return nil, err
	}

	return mentions, nil
}

// ClearNotifications - Clear all notifications
func ClearNotifications(instanceURL, accessToken string) error {
	u, err := url.Parse(instanceURL + "/api/v1/notifications/clear")
	if err != nil {
		return err
	}

	_, err = resty.New().R().
		SetAuthToken(accessToken).
		Get(u.String())
	if err != nil {
		return err
	}

	return nil
}
