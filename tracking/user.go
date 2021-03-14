package tracking

import (
	"github.com/eltrufas/casino/rewards"
	"log"
	"time"
)

type userTracker struct {
	k            userKey
	updates      chan bool
	done         chan userKey
	rewards      chan rewards.Reward
	interval     time.Duration
	rewardAmount int
}

func newUserTracker(k userKey, interval time.Duration, rewardAmount int, done chan userKey, rewards chan rewards.Reward) userTracker {
	updates := make(chan bool)
	t := userTracker{
		updates:      updates,
		k:            k,
		done:         done,
		interval:     interval,
		rewardAmount: rewardAmount,
		rewards:      rewards,
	}
	return t
}

func (t *userTracker) SendUpdate(u bool) {
	t.updates <- u
}

func (t userTracker) loop() {
	connected := false
	ticker := time.NewTicker(t.interval)
	log.Printf("Starting loop")
	defer ticker.Stop()
	for {
		select {
		case u := <-t.updates:
			log.Printf("Handling update %v", u)
			if u != connected {
				ticker.Reset(t.interval)
			}
			connected = u
		case <-ticker.C:
			if !connected {
				t.done <- t.k
				return
			}
			r := rewards.Reward{
				UserID:  t.k.UserID,
				GuildID: t.k.GuildID,
				Amount:  t.rewardAmount,
			}
			log.Printf("Paying out reward")
			t.rewards <- r
		case <-t.done:
			t.shutdown()
			return
		}
	}
}

func (t userTracker) shutdown() {
	close(t.updates)
}

func (t userTracker) Stop() {
	t.done <- t.k
}
