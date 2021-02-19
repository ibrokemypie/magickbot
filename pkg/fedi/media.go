package fedi

import (
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/go-resty/resty/v2"
)

// Attachment - Mastodon attachment object
type Attachment struct {
	ID        string `json:"id"`
	MediaType string `json:"type"`
	URL       string `json:"url"`
	RemoteURL string `json:"remote_url"`
}

// GetMedia - Download media to tmp from url
func GetMedia(mediaURL string, accessToken string) (string, error) {
	filename := "/tmp/" + path.Base(mediaURL)

	_, err := resty.New().R().
		SetAuthToken(accessToken).
		SetOutput(filename).
		Get(mediaURL)
	if err != nil {
		return "", err
	}

	return filename, nil
}

// PostMedia - Upload files and create a new status
func PostMedia(content string, files []string, reply Status, instanceURL, accessToken string) error {

	var mediaIDs = make([]string, 0)

	for _, file := range files {
		u, err := url.Parse(instanceURL + "/api/v1/media")
		if err != nil {
			return err
		}

		var result Attachment

		_, err = resty.New().R().
			SetAuthToken(accessToken).
			SetFile("file", file).
			SetResult(&result).
			Post(u.String())
		if err != nil {
			return err
		}

		mediaIDs = append(mediaIDs, result.ID)
	}

	if len(mediaIDs) > 0 {
		err := PostStatus(content, mediaIDs, reply, instanceURL, accessToken)
		if err != nil {
			return err
		}

		for _, file := range files {
			os.Remove("/tmp/" + strings.TrimPrefix(file, "/tmp/out."))
			os.Remove(file)
		}
	}

	return nil
}
