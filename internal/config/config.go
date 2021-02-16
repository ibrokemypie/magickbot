package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ibrokemypie/magickbot/pkg/auth"
	"github.com/spf13/viper"
)

func LoadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config/magickbot")
	viper.AddConfigPath(".")

	viper.SetDefault("instance.visibility", "public")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config file not found,.creating one at $HOME/.config/magickbot/config.yaml")

			// get the user's home dir
			home, err := os.UserHomeDir()
			if err != nil {
				panic(err)
			}

			// path to config file
			configDir := filepath.Join(home, ".config/magickbot")

			// attempt to create containing folders
			err = os.MkdirAll(configDir, 0777)
			if err != nil {
				panic(err)
			}

			// attempt to write a new config file
			if err := viper.SafeWriteConfig(); err != nil {
				panic(err)
			}
		} else {
			// Config file was found but another error was produced
			panic(err)
		}
	}

	if !viper.IsSet("instance.instance_url") || !viper.IsSet("instance.access_token") {
		instanceURL, accessToken := auth.Authorize()

		viper.Set("instance.instance_url", instanceURL)
		viper.Set("instance.access_token", accessToken)

		viper.WriteConfig()
	}
}
