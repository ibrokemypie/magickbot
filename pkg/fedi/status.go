package fedi

import (
	"net/url"

	"github.com/go-resty/resty/v2"
)

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
		panic(err)
	}

	return result, nil
}

// PostStatus - Posts a text status
func PostStatus(contents, replyToID, instanceURL, accessToken string) error {
	u, err := url.Parse(instanceURL + "/api/v1/statuses")
	if err != nil {
		return (err)
	}

	_, err = resty.New().R().
		SetAuthToken(accessToken).
		SetFormData(map[string]string{
			"in_reply_to_id": replyToID,
			"status":         contents,
		}).
		Post(u.String())
	if err != nil {
		return (err)
	}

	return nil
}
