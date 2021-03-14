package main

import (
	"github.com/eltrufas/casino/cmd"
	"github.com/spf13/viper"
	"log"
)

func main() {
	viper.SetConfigName("casino")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("CASINO_")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Unable to read config file: %", err)
	}

	rootCmd := cmd.BuildRootCommand()
	discordCmd := cmd.BuildDiscordCommand()
	rootCmd.AddCommand(discordCmd)
	rootCmd.Execute()
}
