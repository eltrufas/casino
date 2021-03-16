package discord

import (
	"context"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
)

func (c *Client) GetBalance(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate) {
	userID := m.Author.ID
	guildID := m.GuildID
	balance, err := c.repo.GetUserBalance(ctx, guildID, userID)
	if err != nil {
		log.Printf("Repository.GetUserBalance: %v")
	}
	msg := fmt.Sprintf("<@%s>: %d puntitos", userID, balance.Balance)
	s.ChannelMessageSend(m.ChannelID, msg)
}
