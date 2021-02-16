package bot

import (
	"time"

	"github.com/ibrokemypie/magickbot/pkg/fedi"
	"github.com/spf13/viper"
)

func BotLoop() {
	instanceURL := viper.GetString("instance.instance_url")
	accessToken := viper.GetString("instance.access_token")

	for range time.Tick(time.Second * 5) {
		fedi.GetMentions(instanceURL, accessToken)
	}
}
