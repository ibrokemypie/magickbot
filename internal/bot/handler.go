package bot

import (
	"log"
	"math/rand"
	"strconv"
	"strings"

	"github.com/ibrokemypie/magickbot/pkg/fedi"
	"github.com/ibrokemypie/magickbot/pkg/magick"
	"github.com/microcosm-cc/bluemonday"
	"github.com/spf13/viper"
)

func handleMention(mention fedi.Notification, selfID string, instanceURL, accessToken string) {
	var status fedi.Status
	var operation string
	var argument = 0
	var providedMedia = false
	var maxIterations = viper.GetInt("max_iterations")

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
		p := bluemonday.StrictPolicy().AddSpaceWhenStrippingTag(true)
		text := p.Sanitize(mention.Status.Content)

		textSplit := strings.Fields(text)

		for k, v := range textSplit {
			if v == "help" {
				PostHelp(mention.Status, selfID, instanceURL, accessToken)
				return
			}

			for _, command := range magick.MagickCommands {
				if v == command {
					operation = v
					break
				}
			}

			if v == "random" {
				operation = magick.MagickCommands[rand.Intn(len(magick.MagickCommands))]
				argument = rand.Intn(maxIterations) + 1
			}

			// If the next text is a number, and number is between 1 and 15 inclusive, run this many iterations of command
			if len(textSplit) > k+1 {
				j, err := strconv.Atoi(textSplit[k+1])
				if err == nil {
					argument = j
				}
			}

			if operation != "" {
				break
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

		// Try to run the magick operation on the files
		argument, err := magick.RunMagick(operation, files, argument)
		// retry once
		if err != nil {
			log.Println(err)
			argument, err = magick.RunMagick(operation, files, argument)
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
		if argument != -1 {
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
	}
}
