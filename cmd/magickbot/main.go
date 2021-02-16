package main

import (
	"github.com/ibrokemypie/magickbot/internal/bot"
	"github.com/ibrokemypie/magickbot/internal/config"
)

func main() {
	config.LoadConfig()
	bot.BotLoop()
}
