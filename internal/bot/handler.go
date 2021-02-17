package bot

import (
	"log"
	"strconv"
	"strings"

	"github.com/ibrokemypie/magickbot/pkg/fedi"
	"github.com/ibrokemypie/magickbot/pkg/magick"
	"github.com/microcosm-cc/bluemonday"
)

func handleMention(mention fedi.Mention, instanceURL, accessToken string) {
	var status fedi.Status
	var operation func([]string, int) error
	var iterations = 1

	if mention.Status.MediaAttachments != nil && len(mention.Status.MediaAttachments) > 0 {
		status = mention.Status
	} else {
		reply, err := fedi.GetStatus(mention.Status.ReplyToID, instanceURL, accessToken)
		if err != nil {
			PostError(err, mention.Status.ID, instanceURL, accessToken)
			return
		}

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
				operation = magick.Explode
			case "implode":
				operation = magick.Implode
			}

			if operation != nil {
				// If the next text is a number, and number is between 1 and 15 inclusive, run this many iterations of command
				if len(textSplit) > k+1 {
					j, err := strconv.Atoi(textSplit[k+1])
					if err == nil && iterations > 0 && iterations <= 15 {
						iterations = j
					}
				}

				// For each attached media, download it and add to files list, then run the command on the files list, finally posting the files in a reply
				for _, attachment := range status.MediaAttachments {
					files = append(files, fedi.GetMedia(attachment.URL, accessToken))
					err := operation(files, iterations)
					// retry once
					if err != nil {
						log.Println(err)
						err = operation(files, iterations)
						if err != nil {
							PostError(err, mention.Status.ID, instanceURL, accessToken)
							return
						}
					}

					err = fedi.PostMedia(files, mention.Status.ID, instanceURL, accessToken)
					if err != nil {
						PostError(err, mention.Status.ID, instanceURL, accessToken)
						return
					}
				}
				break Loop
			}
		}
	}
}
