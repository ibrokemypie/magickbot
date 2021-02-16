package bot

import (
	"fmt"
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

		for _, mention := range mentions {
			for _, media := range mention.Status.MediaAttachments {
				fmt.Println(media.URL)
			}
		}
	}
}
