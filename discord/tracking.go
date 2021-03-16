package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/eltrufas/casino/tracking"
	"log"
)

func isActive(u *discordgo.VoiceStateUpdate) bool {
	connected := u.ChannelID != ""

	return connected && !u.Mute && !u.Deaf && !u.SelfMute && !u.SelfDeaf

}

func (c *Client) HandleVoiceStateUpdate(s *discordgo.Session, m *discordgo.VoiceStateUpdate) {
	log.Printf("Handling event %v", m)

	u := tracking.UserUpdate{
		UserID:    m.UserID,
		GuildID:   m.GuildID,
		Connected: isActive(m),
	}
	c.tracker.Enqueue(u)
}
