package fedi

import (
	"net/url"
	"strconv"

	"github.com/go-resty/resty/v2"
)

// Status - Mastodon status object
type Status struct {
	ID               string       `json:"id"`
	ReplyToID        string       `json:"in_reply_to_id"`
	Content          string       `json:"content"`
	Text             string       `json:"text"`
	MediaAttachments []Attachment `json:"media_attachments"`
	Account          Account      `json:"account"`
	Sensitive        bool         `json:"sensitive"`
	Visibility       string       `json:"visibility"`
	Mentions         []Account    `json:"mentions"`
}

// GetStatus - Return a status object from an ID
func GetStatus(id, instanceURL, accessToken string) (Status, error) {
	url, err := url.Parse(instanceURL + "/api/v1/statuses/" + id)
	if err != nil {
		return Status{}, err
	}

	var result Status

	_, err = resty.New().R().
		SetAuthToken(accessToken).
		SetResult(&result).
		Get(url.String())
	if err != nil {
		return Status{}, err
	}

	return result, nil
}

// PostStatus - Posts a text status
func PostStatus(contents string, reply Status, instanceURL, accessToken string) error {
	u, err := url.Parse(instanceURL + "/api/v1/statuses")
	if err != nil {
		return err
	}

	_, err = resty.New().R().
		SetAuthToken(accessToken).
		SetFormData(map[string]string{
			"in_reply_to_id": reply.ID,
			"status":         contents,
			"visibility":     reply.Visibility,
			"sensitive":      strconv.FormatBool(reply.Sensitive),
		}).
		Post(u.String())
	if err != nil {
		return err
	}

	return nil
}
