package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/eltrufas/casino/tracking"
	"log"
)

type DiscordVoiceTracker struct {
	t tracking.Tracker
}

func NewDiscordVoiceTracker(t tracking.Tracker) *DiscordVoiceTracker {
	return &DiscordVoiceTracker{
		t: t,
	}
}

func isActive(u *discordgo.VoiceStateUpdate) bool {
	connected := u.ChannelID != ""

	return connected && !u.Mute && !u.Deaf && !u.SelfMute && !u.SelfDeaf

}

func (dt DiscordVoiceTracker) HandleVoiceStateUpdate(s *discordgo.Session, m *discordgo.VoiceStateUpdate) {
	log.Printf("Handling event %v", m)

	u := tracking.UserUpdate{
		UserID:    m.UserID,
		GuildID:   m.GuildID,
		Connected: isActive(m),
	}
	dt.t.Enqueue(u)
}
