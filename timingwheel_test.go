package timingwheel_test

import (
	"testing"
	"time"

	"github.com/RussellLuo/timingwheel"
)

func TestTimingWheel_Add(t *testing.T) {
	tw := timingwheel.NewTimingWheel(time.Millisecond, 20)
	tw.Start()
	defer tw.Stop()

	durations := []time.Duration{
		1 * time.Millisecond,
		5 * time.Millisecond,
		10 * time.Millisecond,
		50 * time.Millisecond,
		100 * time.Millisecond,
		500 * time.Millisecond,
		1 * time.Second,
	}
	for _, d := range durations {
		t.Run("", func(t *testing.T) {
			exitC := make(chan time.Time)

			timer := timingwheel.AfterFunc(d, func() {
				exitC <- time.Now()
			})
			tw.Add(timer)

			got := (<-exitC).Truncate(time.Millisecond)
			min := time.Unix(0, timer.Expiration*int64(time.Millisecond))

			err := 5 * time.Millisecond
			if got.Before(min) || got.After(min.Add(err)) {
				t.Errorf("NewTimer(%s) want [%s, %s], got %s", d, min, min.Add(err), got)
			}
		})
	}
}
