package discord

import (
	"context"
	"github.com/bwmarrin/discordgo"
)

type MessageHandler func(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate)

func BuildHandler(h MessageHandler) func(s *discordgo.Session, m *discordgo.MessageCreate) {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		ctx := context.Background()
		h(ctx, s, m)
	}
}

func ComposeMiddleware(hs ...func(MessageHandler) MessageHandler) func(MessageHandler) MessageHandler {
	return func(h MessageHandler) MessageHandler {
		for i := len(hs) - 1; i >= 0; i-- {
			h = hs[i](h)
		}
		return h
	}
}
