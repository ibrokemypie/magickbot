package bot

import (
	"strings"

	"github.com/ibrokemypie/magickbot/pkg/magick"

	"github.com/ibrokemypie/magickbot/pkg/fedi"
)

func handleMention(mention fedi.Mention, instanceURL, accessToken string) {
	if len(mention.Status.MediaAttachments) != 0 {
		files := make([]string, 0)

		if strings.Contains(strings.ToLower(mention.Status.Content), " explode") {
			for _, attachment := range mention.Status.MediaAttachments {
				files = append(files, fedi.GetMedia(attachment.URL, accessToken))
				magick.Explode(files)
				fedi.PostMedia(files, mention.Status.ID, instanceURL, accessToken)
			}

		} else if strings.Contains(strings.ToLower(mention.Status.Content), " implode") {
			for _, attachment := range mention.Status.MediaAttachments {
				files = append(files, fedi.GetMedia(attachment.URL, accessToken))
				magick.Implode(files)
				fedi.PostMedia(files, mention.Status.ID, instanceURL, accessToken)
			}
		}
	}
}
