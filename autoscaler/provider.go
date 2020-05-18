package autoscaler

import "time"

type Provider interface {
	SetCapacity(int) error
	SetMinimumAge(duration time.Duration) error
}
