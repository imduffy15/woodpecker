package autoscaler

import (
	"context"
	"github.com/laszlocph/woodpecker/autoscaler"
	droneserver "github.com/laszlocph/woodpecker/server"
	"math"
	"sync"
	"time"
)

type scaler struct {
	capacityPerAgent int
	provider         autoscaler.Provider
}

func New(capacityPerAgent int, minimumAge time.Duration, provider autoscaler.Provider) autoscaler.Autoscaler {
	provider.SetMinimumAge(minimumAge)
	return &scaler{
		capacityPerAgent: capacityPerAgent,
		provider:         provider,
	}
}

func (a scaler) Start(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		a.scale(ctx)
		wg.Done()
	}()
	wg.Wait()
}

func (a scaler) scale(ctx context.Context) {
	const interval = time.Second * 10
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(interval):
			queue := droneserver.Config.Services.Queue.Info(nil).Stats
			pending := queue.Pending
			running := queue.Running
			capacity := queue.Workers * a.capacityPerAgent

			free := autoscaler.Max(capacity-running, 0)
			diff := int(math.Ceil(float64(pending-free) / float64(a.capacityPerAgent)))

			var desired = queue.Workers

			if diff > 0 {
				desired = queue.Workers + diff
			}

			if diff < 0 {
				desired = queue.Workers - autoscaler.Abs(diff)
			}

			a.provider.SetCapacity(desired)
		}
	}
}
