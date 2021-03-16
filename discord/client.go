package discord

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/eltrufas/casino/models"
	"github.com/eltrufas/casino/tracking"
)

type Client struct {
	repo    models.Repository
	tracker tracking.Tracker
}

func NewClient(repo models.Repository, tracker tracking.Tracker) Client {
	return Client{
		repo:    repo,
		tracker: tracker,
	}
}

func (c *Client) HandleMessage(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate) {
	// TODO: Actual command parsing
	if m.Content == "bal" || m.Content == "balance" {
		c.GetBalance(ctx, s, m)
	}
}
