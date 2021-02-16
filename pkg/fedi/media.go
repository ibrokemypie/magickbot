package fedi

import (
	"log"
	"net/url"
	"os"
	"path"

	"github.com/go-resty/resty/v2"
)

type Media struct {
	ID string `json:"id"`
}

// GetMedia - Download media to tmp from url
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

// PostMedia - Upload files and create a new status
func PostMedia(files []string, replyToID, instanceURL, accessToken string) error {

	var mediaIDs = make([]string, 0)

	for _, file := range files {
		u, err := url.Parse(instanceURL + "/api/v1/media")
		if err != nil {
			panic(err)
		}

		filename := "/tmp/out." + path.Base(file)

		var result Media

		_, err = resty.New().R().
			SetAuthToken(accessToken).
			SetFile("file", filename).
			SetResult(&result).
			Post(u.String())
		if err != nil {
			log.Fatal(err)
			return (err)
		}

		mediaIDs = append(mediaIDs, result.ID)

		os.Remove("/tmp/" + path.Base(file))
		os.Remove(filename)
	}

	if len(mediaIDs) > 0 {
		u, err := url.Parse(instanceURL + "/api/v1/statuses")
		if err != nil {
			panic(err)
		}

		_, err = resty.New().R().
			SetAuthToken(accessToken).
			SetFormDataFromValues(url.Values{
				"in_reply_to_id": []string{replyToID},
				"media_ids[]":    mediaIDs,
			}).
			Post(u.String())
	}

	return nil
}
