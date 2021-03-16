package discord

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"log"
)

func LogMessage(next MessageHandler) MessageHandler {
	return func(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate) {
		log.Printf(m.Content)
		next(ctx, s, m)
	}
}
