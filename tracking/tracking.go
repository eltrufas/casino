package tracking

import (
	"github.com/eltrufas/casino/rewards"
	"log"
	"time"
)

type Tracker interface {
	Start()
	Enqueue(u UserUpdate)
	RewardChannel() chan rewards.Reward
	Stop()
}

type tracker struct {
	users        map[userKey]userTracker
	interval     time.Duration
	inbox        chan UserUpdate
	clear        chan userKey
	rewards      chan rewards.Reward
	done         chan struct{}
	rewardAmount int
}

func New(interval time.Duration, rewardAmount int) (Tracker, error) {
	users := make(map[userKey]userTracker)
	inbox := make(chan UserUpdate, 512)
	clear := make(chan userKey)
	rewards := make(chan rewards.Reward, 1024)
	done := make(chan struct{})
	return &tracker{
		users:        users,
		interval:     interval,
		rewardAmount: rewardAmount,
		inbox:        inbox,
		clear:        clear,
		rewards:      rewards,
		done:         done,
	}, nil
}

func (t *tracker) Enqueue(u UserUpdate) {
	t.inbox <- u
}

func (t *tracker) RewardChannel() chan rewards.Reward {
	return t.rewards
}

func (t *tracker) Start() {
	go t.loop()
}

func (t *tracker) loop() {
	for {
		select {
		case u := <-t.inbox:
			log.Printf("Processing update %v", u)
			t.handleUpdate(u)
		case k := <-t.clear:
			log.Printf("Clearing %v", k)
			t.handleClear(k)
		case <-t.done:
			log.Printf("Shutting down")
			t.handleShutdown()
			return
		}
	}
}

func (t *tracker) Stop() {
	t.done <- struct{}{}
}

func (t *tracker) handleUpdate(e UserUpdate) {
	// This ain't threadsafe. Should only be called from the loop
	k := e.getKey()
	u, ok := t.users[k]
	if !ok {
		if !e.Connected {
			// Don't do anything for a disconnect user that we're not tracking
			return
		}
		u = t.launchUserTracker(k)
		u.SendUpdate(e.Connected)
		return
	}
	u.SendUpdate(e.Connected)
}

func (t *tracker) handleClear(k userKey) {
	delete(t.users, k)
}

func (t *tracker) handleShutdown() {
	for _, v := range t.users {
		v.Stop()
	}
}

func (t *tracker) launchUserTracker(k userKey) userTracker {
	ut := newUserTracker(k, t.interval, t.rewardAmount, t.clear, t.rewards)
	go ut.loop()
	t.users[k] = ut
	return ut
}
