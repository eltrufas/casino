package rewards

import (
	"context"
	"github.com/eltrufas/casino/models"
	"log"
)

type Reward struct {
	UserID  string
	GuildID string
	Amount  int
}

type RewardWorker struct {
	repo  models.Repository
	queue chan Reward
}

func NewRewardWorker(repo models.Repository, queue chan Reward) RewardWorker {
	return RewardWorker{
		repo:  repo,
		queue: queue,
	}
}

func (w RewardWorker) LaunchWorkers(ctx context.Context, workerCount int) {
	for i := 0; i < workerCount; i++ {
		go w.workerLoop(ctx)
	}
}

func (w RewardWorker) workerLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case r := <-w.queue:
			txn := models.Transaction{
				UserID:  r.UserID,
				GuildID: r.GuildID,
				Amount:  int64(r.Amount),
				Note:    "Voice channel reward",
			}
			err := w.repo.CreateTransaction(ctx, &txn)
			if err != nil {
				log.Printf("CreateTransaction: %v", err)
			}
			log.Printf("Created txn %v", txn)
		}
	}
}
