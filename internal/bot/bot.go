package bot

import (
	"log"
	"strings"
	"time"

	"github.com/ibrokemypie/magickbot/pkg/fedi"
	"github.com/spf13/viper"
)

// BotLoop - main loop of the bot
func BotLoop() {
	instanceURL := viper.GetString("instance.instance_url")
	accessToken := viper.GetString("instance.access_token")

	self, err := fedi.GetCurrentUser(instanceURL, accessToken)
	if err != nil {
		PostError(err, fedi.Status{}, instanceURL, accessToken)
		return
	}

	for range time.Tick(time.Second * 5) {
		mentions := fedi.GetMentions(instanceURL, accessToken)

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

	content.WriteString("Magickbot Help: \n\n")
	content.WriteString("Usage: \n")
	content.WriteString("Tag the bot either in a status containing media, a reply to a status containing media, a status containing mentions of users to apply to their avatars or a reply to a status with no media to apply to the user's avatar. Include the command (eg. explode) in your status, optionally include an argument. The only order that matters is argument must be after command.\n\n")
	content.WriteString("command [argument] [@user...]\n\n")
	content.WriteString("Commands: \n")
	content.WriteString("help\n")
	content.WriteString("explode [iterations]\n")
	content.WriteString("implode [iterations]\n")
	content.WriteString("magick [scale]\n")

	err := fedi.PostStatus(content.String(), reply, instanceURL, accessToken)
	if err != nil {
		PostError(err, reply, instanceURL, accessToken)
	}
}
