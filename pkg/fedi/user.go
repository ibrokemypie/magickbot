package fedi

import (
	"net/url"

	"github.com/go-resty/resty/v2"
)

// Account - Mastodon account object
type Account struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	URL      string `json:"url"`
	Acct     string `json:"acct"`
	Avatar   string `json:"avatar"`
}

// GetUser - Return a User object from an ID
func GetUser(id, instanceURL, accessToken string) (Account, error) {
	url, err := url.Parse(instanceURL + "/api/v1/accounts/" + id)
	if err != nil {
		return Account{}, err
	}

	var result Account

	_, err = resty.New().R().
		SetAuthToken(accessToken).
		SetResult(&result).
		Get(url.String())
	if err != nil {
		panic(err)
	}

	return result, nil
}

// GetCurrentUser - Return the current user's user object
func GetCurrentUser(instanceURL, accessToken string) (Account, error) {
	url, err := url.Parse(instanceURL + "/api/v1/accounts/verify_credentials")
	if err != nil {
		return Account{}, err
	}

	var result Account

	_, err = resty.New().R().
		SetAuthToken(accessToken).
		SetResult(&result).
		Get(url.String())
	if err != nil {
		panic(err)
	}

	return result, nil
}
