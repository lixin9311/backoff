package backoff

import (
	"testing"
	"time"
)

func TestDuration(t *testing.T) {

	if e, g := 1000*time.Millisecond, DefaultExponential.Backoff(0); e != g {
		t.Fatalf("expect %s, got %s", e, g)
	}

	between(t, DefaultExponential.Backoff(1), 1000*time.Millisecond, 2000*time.Millisecond)
	between(t, DefaultExponential.Backoff(2), 100*time.Millisecond, 3500*time.Millisecond)
}

func between(t *testing.T, got, low, high time.Duration) {
	if got < low {
		t.Fatalf("expect >= %s, got %s", low, got)
	}
	if got > high {
		t.Fatalf("expect <= %s, got %s", high, got)
	}
}
