package cmd

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/eltrufas/casino/discord"
	"github.com/eltrufas/casino/models"
	"github.com/eltrufas/casino/rewards"
	"github.com/eltrufas/casino/tracking"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func BuildDiscordCommand() *cobra.Command {
	viper.SetDefault("reward.interval", "1m")
	viper.SetDefault("reward.amount", "10")
	return &cobra.Command{
		Use: "discord",
		Run: func(cmd *cobra.Command, args []string) {
			db, err := sqlx.Connect("sqlite3", "./casino.db")
			if err != nil {
				log.Fatalf("sqlx.Conntect: %v", err)
			}
			repo := models.NewRepository(db)

			tracker, err := tracking.New(viper.GetDuration("reward.interval"), viper.GetInt("reward.amount"))
			if err != nil {
				log.Fatalf("Can't create tracker: %w", err)
			}

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			session, err := discordgo.New(viper.GetString("discord.token"))
			client := discord.NewClient(repo, tracker)

			session.AddHandler(client.HandleVoiceStateUpdate)
			middleware := discord.ComposeMiddleware(
				discord.LogMessage,
				discord.RequireMention,
			)
			h := discord.BuildHandler(middleware(discord.MessageHandler(client.HandleMessage)))
			session.AddHandler(h)
			err = session.Open()
			if err != nil {
				log.Fatalf("Unable to connect to discord: %v", err)
			}
			defer session.Close()

			tracker.Launch(ctx)

			rewardChan := tracker.RewardChannel()
			rewardWorker := rewards.NewRewardWorker(repo, rewardChan)

			rewardWorker.LaunchWorkers(ctx, 1)

			sc := make(chan os.Signal, 1)
			signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
			<-sc
		},
	}
}
