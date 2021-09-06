package backoff

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"
)

var (
	errRetry = errors.New("should retry")
	errFail  = errors.New("should fail")
)

type retrier struct {
	Backoff
	max int
}

func (r *retrier) Retry(n int, err error) (time.Duration, bool) {
	if n >= r.max {
		return 0, false
	} else if err != errRetry {
		return 0, false
	}
	d := r.Backoff.Backoff(n)
	log.Printf("%d-th attepmt, retry in %v", n, d)
	return d, true
}

type api struct {
	errs  []error
	index int
}

func (a *api) Call(ctx context.Context) error {
	err := a.errs[a.index]
	a.index++
	return err
}

func TestBackoff(t *testing.T) {
	api := &api{errs: []error{
		errRetry,
		errRetry,
		errRetry,
		nil,
	}}
	retrier := &retrier{
		max: 3,
	}
	err := Invoke(context.Background(), api.Call, retrier.Retry)
	if err != nil {
		t.Error("expect nil error")
	}
}

func TestBackoffFail(t *testing.T) {
	api := &api{errs: []error{
		errRetry,
		errRetry,
		errFail,
		nil,
	}}

	retrier := &retrier{
		max: 3,
	}
	err := Invoke(context.Background(), api.Call, retrier.Retry)
	if err != errFail {
		t.Error("expect failed error")
	}
}

func TestBackoffTimeout(t *testing.T) {
	api := &api{errs: []error{
		errRetry,
		errRetry,
		errRetry,
		nil,
	}}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	retrier := &retrier{
		max: 3,
	}
	err := Invoke(ctx, api.Call, retrier.Retry)
	if err != context.DeadlineExceeded {
		t.Error("expect DeadlineExceeded error")
	}
}
