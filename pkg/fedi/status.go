package fedi

import (
	"net/url"

	"github.com/go-resty/resty/v2"
)

// GetStatus - Return a status object from an ID
func GetStatus(id, instanceURL, accessToken string) Status {
	url, err := url.Parse(instanceURL + "/api/v1/statuses/" + id)
	if err != nil {
		panic(err)
	}

	var result Status

	_, err = resty.New().R().
		SetAuthToken(accessToken).
		SetResult(&result).
		Get(url.String())
	if err != nil {
		panic(err)
	}

	return result
}
