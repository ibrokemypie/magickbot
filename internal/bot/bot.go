package bot

import (
	"time"

	"github.com/ibrokemypie/magickbot/pkg/fedi"
	"github.com/spf13/viper"
)

// BotLoop - main loop of the bot
func BotLoop() {
	instanceURL := viper.GetString("instance.instance_url")
	accessToken := viper.GetString("instance.access_token")

	for range time.Tick(time.Second * 5) {
		mentions := fedi.GetMentions(instanceURL, accessToken)

		for k := range mentions {
			mention := mentions[len(mentions)-1-k]
			handleMention(mention, instanceURL, accessToken)
			viper.Set("last_mention_id", mention.ID)
			viper.WriteConfig()
		}
	}
}
