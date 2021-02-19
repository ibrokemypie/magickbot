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
func PostStatus(contents string, mediaIDs []string, reply Status, instanceURL, accessToken string) error {
	u, err := url.Parse(instanceURL + "/api/v1/statuses")
	if err != nil {
		return err
	}

	_, err = resty.New().R().
		SetAuthToken(accessToken).
		SetFormDataFromValues(url.Values{
			"in_reply_to_id": []string{reply.ID},
			"status":         []string{contents},
			"visibility":     []string{reply.Visibility},
			"sensitive":      []string{strconv.FormatBool(reply.Sensitive)},
			"media_ids[]":    mediaIDs,
		}).
		Post(u.String())
	if err != nil {
		return err
	}

	return nil
}
