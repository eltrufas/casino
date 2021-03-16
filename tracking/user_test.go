package tracking

import (
	"github.com/eltrufas/casino/rewards"
	"testing"
	"time"
)

func buildUserTracker(t *testing.T, configure func(*userTrackerConfig)) userTracker {
	interval, err := time.ParseDuration("20ms")
	if err != nil {
		panic("what")
	}
	config := userTrackerConfig{
		interval:     interval,
		rewardAmount: 10,
		key: userKey{
			UserID:  "123",
			GuildID: "456",
		},

		clear:   make(chan userKey, 1),
		rewards: make(chan rewards.Reward),
	}

	if configure != nil {
		configure(&config)
	}
	return newUserTracker(config)
}

func TestSendUpdate(t *testing.T) {
	timeout, _ := time.ParseDuration("5s")
	updateHookCalled := make(chan struct{})
	tracker := buildUserTracker(t, func(c *userTrackerConfig) {
		c.updateHook = func(u bool) {
			updateHookCalled <- struct{}{}
		}
	})
	go tracker.loop()

	tracker.SendUpdate(false)
	select {
	case <-time.After(timeout):
		t.Error("Update hook not called")
	case <-updateHookCalled:
	}

	if len(updateHookCalled) != 0 {
		t.Error("Update hook called more than once")
	}
	tracker.Stop()

	// shuts down when interval passes and user is disconnected
	doneHookCalled := make(chan struct{})
	tracker = buildUserTracker(t, func(c *userTrackerConfig) {
		c.doneHook = func() {
			doneHookCalled <- struct{}{}
		}
	})

	go tracker.loop()

	tracker.SendUpdate(false)
	select {
	case <-time.After(timeout):
		t.Error("Done hook not called")
	case <-doneHookCalled:
	}

	rewardChan := make(chan rewards.Reward)
	tracker = buildUserTracker(t, func(c *userTrackerConfig) {
		c.rewards = rewardChan
	})

	go tracker.loop()

	tracker.SendUpdate(true)

	expected := rewards.Reward{
		UserID:  "123",
		GuildID: "456",
		Amount:  10,
	}
	select {
	case <-time.After(timeout):
		t.Error("Done hook not called")
	case r := <-rewardChan:
		if r != expected {
			t.Errorf("Expected %#v but got %#v", expected, r)
		}
	}
}
