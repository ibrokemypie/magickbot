package bot

import (
	"log"
	"strconv"
	"strings"

	"github.com/ibrokemypie/magickbot/pkg/fedi"
	"github.com/ibrokemypie/magickbot/pkg/magick"
	"github.com/microcosm-cc/bluemonday"
)

func handleMention(mention fedi.Notification, selfID string, instanceURL, accessToken string) {
	var status fedi.Status
	var operation magick.MagickCommand
	var argument = 0
	var providedMedia = false

	status = mention.Status

	// If mentioning status has no images and reply exists, use reply
	if (mention.Status.MediaAttachments != nil || len(mention.Status.MediaAttachments) > 0) && mention.Status.ReplyToID != "" {
		reply, err := fedi.GetStatus(mention.Status.ReplyToID, instanceURL, accessToken)
		if err != nil {
			PostError(err, mention.Status, instanceURL, accessToken)
			return
		}

		if reply.ID != "" {
			status = reply
		}
		// otherwise apply to the profile pictures of tagged users
	} else if len(mention.Status.Mentions) > 1 {
		// add the profile pics of non self mentioned users as attachments to the status
		for _, m := range mention.Status.Mentions {
			if m.ID != selfID {
				user, err := fedi.GetUser(m.ID, instanceURL, accessToken)
				if err != nil {
					PostError(err, mention.Status, instanceURL, accessToken)
					return
				}

				newImage := fedi.Attachment{URL: user.Avatar}
				status.MediaAttachments = append(status.MediaAttachments, newImage)
			}
		}
	}

	if status.MediaAttachments != nil && len(status.MediaAttachments) > 0 {
		providedMedia = true
	}

	if status.MediaAttachments != nil {
		files := make([]string, 0)
		p := bluemonday.StrictPolicy()
		text := p.Sanitize(mention.Status.Content)

		textSplit := strings.Split(text, " ")

	Loop:
		for k, v := range textSplit {
			switch v {
			case "help":
				PostHelp(mention.Status, selfID, instanceURL, accessToken)
				return
			case magick.EXPLODE:
				operation = magick.EXPLODE
			case magick.IMPLODE:
				operation = magick.IMPLODE
			case magick.MAGICK:
				operation = magick.MAGICK
			case magick.COMPRESS:
				operation = magick.COMPRESS
			default:
				continue
			}

			if operation != "" {
				// If the next text is a number, and number is between 1 and 15 inclusive, run this many iterations of command
				if len(textSplit) > k+1 {
					j, err := strconv.Atoi(textSplit[k+1])
					if err == nil {
						argument = j
					}
				}

				// If there was an attachment in the mention or the status it replied to, use that, otherwise apply operation to the avatar
				if providedMedia {
					// For each attached media, download it and add to files list, then run the command on the files list, finally posting the files in a reply
					for _, attachment := range status.MediaAttachments {
						files = append(files, fedi.GetMedia(attachment.URL, accessToken))
					}
				} else {
					files = append(files, fedi.GetMedia(status.Account.Avatar, accessToken))
				}
			}

			// Try to run the magick operation on the files
			err := magick.RunMagick(operation, files, argument)
			// retry once
			if err != nil {
				log.Println(err)
				err = magick.RunMagick(operation, files, argument)
				if err != nil {
					PostError(err, mention.Status, instanceURL, accessToken)
					return
				}
			}

			content := strings.Builder{}
			for _, m := range mention.Status.Mentions {
				if m.ID != selfID && m.ID != mention.Status.Account.Acct {
					content.WriteString("@")
					content.WriteString(m.Acct)
					content.WriteString(", ")
				}
			}

			content.WriteString("@")
			content.WriteString(mention.Status.Account.Acct)
			content.WriteString("\n")

			content.WriteString("Ran ")
			content.WriteString(string(operation))
			if argument != 0 {
				content.WriteString(" ")
				content.WriteString(strconv.Itoa(argument))
			}
			content.WriteString(":")

			// Try to post the manipulated files
			err = fedi.PostMedia(content.String(), files, mention.Status, instanceURL, accessToken)
			if err != nil {
				PostError(err, mention.Status, instanceURL, accessToken)
				return
			}

			break Loop
		}
	}
}
