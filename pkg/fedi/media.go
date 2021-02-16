package fedi

import (
	"path"

	"github.com/go-resty/resty/v2"
)

func GetMedia(mediaURL string, accessToken string) string {
	filename := "/tmp/" + path.Base(mediaURL)

	_, err := resty.New().R().
		SetAuthToken(accessToken).
		SetOutput(filename).
		Get(mediaURL)
	if err != nil {
		panic(err)
	}

	return filename
}
