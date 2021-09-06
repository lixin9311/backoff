package backoff

import (
	"context"
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Retryer is used by Invoke to determine retry behavior.
// It should report whether a request should be retriedand how long to pause before retrying
// if the previous attempt returned with err. Invoke never calls Retry with nil error.
type Retryer func(n int, err error) (pause time.Duration, shouldRetry bool)

// Sleep is similar to time.Sleep, but it can be interrupted by ctx.Done() closing.
// If interrupted, Sleep returns ctx.Err().
func Sleep(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	select {
	case <-ctx.Done():
		t.Stop()
		return ctx.Err()
	case <-t.C:
		return nil
	}
}

// Backoff implements exponential backoff.
// The wait time between retries is a random value between 0 and the "retry envelope".
// The envelope starts at Initial and increases by the factor of Multiplier every retry,
// but is capped at Max.
type Backoff struct {
	// BaseDelay is the amount of time to backoff after the first failure.
	BaseDelay time.Duration
	// Multiplier is the factor with which to multiply backoffs after a
	// failed retry. Should ideally be greater than 1.
	Multiplier float64
	// Jitter is the factor with which backoffs are randomized.
	Jitter float64
	// MaxDelay is the upper bound of backoff delay.
	MaxDelay time.Duration
}

// Backoff returns the next time.Duration that the caller should use to backoff.
func (bo *Backoff) Backoff(retries int) time.Duration {
	if bo.BaseDelay == 0 {
		bo.BaseDelay = time.Second
	}

	if bo.MaxDelay == 0 {
		bo.MaxDelay = 30 * time.Second
	}
	if bo.Multiplier < 1 {
		bo.Multiplier = 1.6
	}

	if retries == 0 {
		return bo.BaseDelay
	}
	backoff, max := float64(bo.BaseDelay), float64(bo.MaxDelay)
	for backoff < max && retries > 0 {
		backoff *= bo.Multiplier
		retries--
	}
	if backoff > max {
		backoff = max
	}
	// Randomize backoff delays so that if a cluster of requests start at
	// the same time, they won't operate in lockstep.
	backoff *= 1 + bo.Jitter*(rand.Float64()*2-1)
	if backoff < 0 {
		return 0
	}
	return time.Duration(backoff)
}

// invoke implements Invoke, taking an additional sleeper argument for testing.
func Invoke(ctx context.Context, call func(context.Context) error, retryer Retryer) error {
	for n := 0; ; n++ {
		err := call(ctx)
		if err == nil {
			return nil
		}
		if retryer == nil {
			return err
		}
		// Never retry permanent certificate errors. (e.x. if ca-certificates
		// are not installed). We should only make very few, targeted
		// exceptions: many (other) status=Unavailable should be retried, such
		// as if there's a network hiccup, or the internet goes out for a
		// minute. This is also why here we are doing string parsing instead of
		// simply making Unavailable a non-retried code elsewhere.
		if strings.Contains(err.Error(), "x509: certificate signed by unknown authority") {
			return err
		}

		if d, ok := retryer(n, err); !ok {
			return err
		} else if err = Sleep(ctx, d); err != nil {
			return err
		}
	}
}
