package limiter

import (
	"testing"
	"time"
)

func TestLimiter(t *testing.T) {
	l := New(5)
	go func() {
		for i := 0; i < 6; i++ {
			l.Allow()
			t.Log(i, "allow")
		}
	}()
	for i := 0; i < 6; i++ {
		time.Sleep(time.Second)
		l.Done()
		t.Log(i, "done")
	}
	l.Wait()
}
