package discord

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

func RequireMention(next MessageHandler) MessageHandler {
	return func(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate) {
		userID := s.State.User.ID
		mentionString := fmt.Sprintf("<@!%s>", userID)

		if !strings.HasPrefix(m.Content, mentionString) {
			return
		}

		// Remove mention from mentions slice
		for i, mention := range m.Mentions {
			if userID == mention.ID {
				copy(m.Mentions[i:], m.Mentions[i+1:])
				break
			}
		}

		m.Content = strings.TrimSpace(m.Content[len(mentionString):])

		next(ctx, s, m)
	}
}
