package tracking

import (
	"github.com/eltrufas/casino/rewards"
	"time"
)

type userTrackerConfig struct {
	interval     time.Duration
	rewardAmount int

	key     userKey
	clear   chan userKey
	rewards chan rewards.Reward

	updateHook func(u bool)
	doneHook   func()
}

type userTracker struct {
	userTrackerConfig

	updates chan bool
	done    chan struct{}
}

func newUserTracker(config userTrackerConfig) userTracker {
	updates := make(chan bool, 8)
	done := make(chan struct{}, 1)
	t := userTracker{
		userTrackerConfig: config,
		updates:           updates,
		done:              done,
	}
	return t
}

func (t *userTracker) SendUpdate(u bool) {
	t.updates <- u
}

func (t userTracker) loop() {
	connected := false
	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()
	for {
		select {
		case u := <-t.updates:
			if u != connected {
				ticker.Reset(t.interval)
			}
			connected = u
			if t.updateHook != nil {
				t.updateHook(connected)
			}
		case <-ticker.C:
			if !connected {
				t.clear <- t.key
				t.shutdown()
				return
			}
			r := rewards.Reward{
				UserID:  t.key.UserID,
				GuildID: t.key.GuildID,
				Amount:  t.rewardAmount,
			}
			t.rewards <- r
		case <-t.done:
			t.shutdown()
			return
		}
	}
}

func (t userTracker) shutdown() {
	if t.doneHook != nil {
		t.doneHook()
	}
}

func (t userTracker) Stop() {
	t.done <- struct{}{}
}
