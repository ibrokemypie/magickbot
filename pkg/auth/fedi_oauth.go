package auth

import (
	"fmt"
	"net/url"

	"github.com/go-resty/resty/v2"
)

type app struct {
	ClientName   string `json:"name"`
	Website      string `json:"website"`
	RedirectURI  string `json:"redirect_uri"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}
type oauth struct {
	AccessToken string `json:"access_token"`
}

func getFediInstance() string {
	fmt.Println("Paste the fedi instance url (eg: https://mastodon.social) and press enter")

	var instance string
	fmt.Scanln(&instance)

	return instance
}

func getApp(instanceURL string) app {
	clientName := "magickbot"
	website := "https://github.com/ibrokemypie/magickbot"
	redirectURI := "urn:ietf:wg:oauth:2.0:oob"

	url, err := url.Parse(instanceURL + "/api/v1/apps")
	if err != nil {
		panic(err)
	}

	var result app

	_, err = resty.New().R().
		SetFormData(map[string]string{
			"client_name":   clientName,
			"scopes":        "read write",
			"website":       website,
			"redirect_uris": redirectURI,
		}).
		SetResult(&result).
		Post(url.String())
	if err != nil {
		panic(err)
	}

	return result
}

func authorizeUser(instanceURL string, app app) string {
	url, err := url.Parse(instanceURL + "/oauth/authorize")
	if err != nil {
		panic(err)
	}

	q := url.Query()
	q.Set("client_id", app.ClientID)
	q.Set("redirect_uri", app.RedirectURI)
	q.Set("response_type", "code")
	q.Set("force_login", "true")
	q.Set("scope", "read write")
	url.RawQuery = q.Encode()

	fmt.Println("Please open the following URL in your browser.")
	fmt.Println("Once you have authenticated, paste the token from the page into this window.")
	fmt.Println(url)

	var token string
	fmt.Scanln(&token)

	return token
}

func oauthToken(instanceURL string, app app, token string) string {
	url, err := url.Parse(instanceURL + "/oauth/token")
	if err != nil {
		panic(err)
	}

	var result map[string]interface{}

	_, err = resty.New().R().
		SetFormData(map[string]string{
			"client_id":     app.ClientID,
			"client_secret": app.ClientSecret,
			"redirect_uri":  app.RedirectURI,
			"scope":         "read write",
			"grant_type":    "authorization_code",
			"code":          token,
		}).
		SetResult(&result).
		Post(url.String())
	if err != nil {
		panic(err)
	}

	if result["access_token"] == nil {
		panic(result)
	}

	return result["access_token"].(string)
}
