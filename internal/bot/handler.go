package bot

import (
	"strconv"
	"strings"

	"github.com/ibrokemypie/magickbot/pkg/fedi"
	"github.com/ibrokemypie/magickbot/pkg/magick"
	"github.com/microcosm-cc/bluemonday"
)

func handleMention(mention fedi.Mention, instanceURL, accessToken string) {
	var status fedi.Status

	if mention.Status.MediaAttachments != nil && len(mention.Status.MediaAttachments) > 0 {
		status = mention.Status
	} else {
		reply := fedi.GetStatus(mention.Status.ReplyToID, instanceURL, accessToken)
		if reply.MediaAttachments != nil && len(reply.MediaAttachments) > 0 {
			status = reply
		}
	}

	if status.MediaAttachments != nil {
		files := make([]string, 0)
		p := bluemonday.StrictPolicy()
		text := p.Sanitize(mention.Status.Content)

		textSplit := strings.Split(text, " ")
	Loop:
		for k, v := range textSplit {
			switch v {
			case "explode":
				// If the next text is a number, and number is between 1 and 15 inclusive, run this many iterations of command
				i := 1
				if len(textSplit) > k+1 {
					j, err := strconv.Atoi(textSplit[k+1])
					if err == nil && i > 0 && i <= 15 {
						i = j
					}
				}

				// For each attached media, download it and add to files list, then run the command on the files list, finally posting the files in a reply
				for _, attachment := range status.MediaAttachments {
					files = append(files, fedi.GetMedia(attachment.URL, accessToken))
					magick.Explode(files, i)
					fedi.PostMedia(files, status.ID, instanceURL, accessToken)
				}
				break Loop
			case "implode":
				// If the next text is a number, and number is between 1 and 15 inclusive, run this many iterations of command
				i := 1
				if len(textSplit) > k+1 {
					j, err := strconv.Atoi(textSplit[k+1])
					if err == nil && i > 0 && i <= 15 {
						i = j
					}
				}

				// For each attached media, download it and add to files list, then run the command on the files list, finally posting the files in a reply
				for _, attachment := range status.MediaAttachments {
					files = append(files, fedi.GetMedia(attachment.URL, accessToken))
				}

				magick.Implode(files, i)
				fedi.PostMedia(files, status.ID, instanceURL, accessToken)
				break Loop
			}
		}
	}
}
