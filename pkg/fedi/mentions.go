package fedi

import (
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"
)

// Account - Mastodon account object
type Account struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Acct     string `json:"acct"`
}

// Attachment - Mastodon attachment object
type Attachment struct {
	ID        string `json:"id"`
	MediaType string `json:"type"`
	URL       string `json:"url"`
	RemoteURL string `json:"remote_url"`
}

// Status - Mastodon status object
type Status struct {
	ID               string       `json:"id"`
	ReplyToID        string       `json:"in_reply_to_id"`
	Content          string       `json:"content"`
	Text             string       `json:"text"`
	MediaAttachments []Attachment `json:"media_attachments"`
}

// Mention - Mastodon mention object
type Mention struct {
	ID      string  `json:"id"`
	Account Account `json:"account"`
	Status  Status  `json:"status"`
}

// GetMentions -
func GetMentions(instanceURL, accessToken string) []Mention {
	fmt.Println("Getting mentions")

	url, err := url.Parse(instanceURL + "/api/v1/notifications")
	if err != nil {
		panic(err)
	}

	lastID := viper.GetString("last_mention_id")
	mentions := make([]Mention, 0)

	_, err = resty.New().R().
		SetAuthToken(accessToken).
		SetFormData(map[string]string{
			"exclude_types": "follow favourite reblog poll follow_request",
		}).
		SetQueryParam("min_id", lastID).
		SetResult(&mentions).
		Get(url.String())
	if err != nil {
		panic(err)
	}
	return mentions
}
