package bot

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/ibrokemypie/magickbot/pkg/fedi"
	"github.com/spf13/viper"
)

// BotLoop - main loop of the bot
func BotLoop() {
	instanceURL := viper.GetString("instance.instance_url")
	accessToken := viper.GetString("instance.access_token")

	for range time.Tick(time.Second * 5) {
		// first check the server is still accessible (and get bot's ID)
		self, err := fedi.GetCurrentUser(instanceURL, accessToken)
		if err != nil {
			PostError(err, fedi.Status{}, instanceURL, accessToken)
			continue
		}

		// next try to get mention notifications
		mentions, err := fedi.GetMentions(instanceURL, accessToken)
		if err != nil {
			PostError(err, fedi.Status{}, instanceURL, accessToken)
			continue
		}

		// if there's no mentions left, clean the notifications
		if len(mentions) == 0 {
			fedi.ClearNotifications(instanceURL, accessToken)
		}

		// handle each mention
		for k := range mentions {
			mention := mentions[len(mentions)-1-k]
			if mention.Account.ID != self.ID {
				handleMention(mention, self.ID, instanceURL, accessToken)
			}

			viper.Set("last_mention_id", mention.ID)
			viper.WriteConfig()
		}
	}
}

// PostError - Helper function to post errors to fedi
func PostError(err error, reply fedi.Status, instanceURL, accessToken string) {
	log.Println(err)
	err = fedi.PostStatus("Magickbot error: "+err.Error(), reply, instanceURL, accessToken)
	if err != nil {
		log.Println(err)
	}
}

// PostHelp - Helper function to post the bot's help
func PostHelp(reply fedi.Status, selfID string, instanceURL, accessToken string) {
	content := strings.Builder{}
	for _, m := range reply.Mentions {
		if m.ID != selfID && m.ID != reply.Account.Acct {
			content.WriteString("@")
			content.WriteString(m.Acct)
			content.WriteString(", ")
		}
	}

	content.WriteString("@")
	content.WriteString(reply.Account.Acct)
	content.WriteString("\n")

	content.WriteString("Magickbot \n\n")
	content.WriteString("Source: \n")
	content.WriteString("https://github.com/ibrokemypie/magickbot \n\n")
	content.WriteString("Usage: \n")
	content.WriteString("Tag the bot either in a status containing media, a reply to a status containing media, a status containing mentions of users to apply to their avatars or a reply to a status with no media to apply to the user's avatar. Include the command (eg. explode) in your status, optionally include an argument. The only order that matters is argument must be after command.\n\n")
	content.WriteString("command [argument] [@user...]\n\n")
	content.WriteString("Commands: \n")
	content.WriteString("help\n")
	content.WriteString("explode [iterations]\n")
	content.WriteString("implode [iterations]\n")
	content.WriteString("magick [scale]\n")
	content.WriteString("moarjpeg [iterations]\n")
	content.WriteString("deepfry\n")
	content.WriteString("random [argument]\n\n")
	content.WriteString("Bot Configuration: \n")
	content.WriteString("Max iterations: ")
	content.WriteString(strconv.Itoa(viper.GetInt("max_iterations")))
	content.WriteString("\nMax input pixels: ")
	content.WriteString(strconv.Itoa(viper.GetInt("max_pixels_in")))
	content.WriteString("\nMax output pixels: ")
	content.WriteString(strconv.Itoa(viper.GetInt("max_pixels_out")))

	err := fedi.PostStatus(content.String(), reply, instanceURL, accessToken)
	if err != nil {
		PostError(err, reply, instanceURL, accessToken)
	}
}
