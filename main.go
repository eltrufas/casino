package main

import (
	"github.com/eltrufas/casino/cmd"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"strings"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	replacer := strings.NewReplacer("-", "_", ".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("Unable to read config file: %", err)
		}
	}

	rootCmd := cmd.BuildRootCommand()
	discordCmd := cmd.BuildDiscordCommand()
	rootCmd.AddCommand(discordCmd)
	rootCmd.Execute()
}
