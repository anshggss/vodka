package mixers

import (
	"testing"
	"time"
)

func TestNewRateLimiter(t *testing.T) {
	vrl := NewRateLimiter(10, 5)
	if vrl == nil {
		t.Fatal("expected non-nil VodkaRateLimiter")
	}
}

func TestRateLimiterAllowsUnderBurst(t *testing.T) {
	vrl := NewRateLimiter(10, 3)

	l := vrl.getVisitor("127.0.0.1")
	for i := 0; i < 3; i++ {
		if !l.allow() {
			t.Fatalf("expected request %d to be allowed within burst", i+1)
		}
	}
}

func TestRateLimiterBlocksOverBurst(t *testing.T) {
	vrl := NewRateLimiter(1, 2)

	l := vrl.getVisitor("10.0.0.1")
	l.allow()
	l.allow()
	if l.allow() {
		t.Fatal("expected request to be blocked after burst exhausted")
	}
}

func TestLazyCleanupRemovesStaleVisitors(t *testing.T) {
	vrl := NewRateLimiter(10, 5)

	vrl.mu.Lock()
	vrl.visitors["stale-ip"] = &visitor{
		l:        newLimiter(10, 5),
		lastSeen: time.Now().Add(-10 * time.Minute),
	}
	vrl.lastCleanup = time.Now().Add(-2 * time.Minute)
	vrl.mu.Unlock()

	vrl.getVisitor("trigger-ip")

	vrl.mu.Lock()
	_, exists := vrl.visitors["stale-ip"]
	vrl.mu.Unlock()

	if exists {
		t.Fatal("expected stale visitor to be cleaned up")
	}
}
