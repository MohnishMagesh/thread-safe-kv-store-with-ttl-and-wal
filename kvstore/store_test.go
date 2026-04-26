package kvstore

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// TestKVStoreRaceCondition hammers the map with concurrent reads and writes
// to prove our sync.RWMutex is working perfectly.
func TestKVStoreRaceCondition(t *testing.T) {
	// 1. Initialize using the constructor so the map and channels are ready,
	// and the background sweeper starts running. We set a long sweep interval
	// so it doesn't interfere with our immediate race test.
	s, _ := NewKVStore(10 * time.Minute)

	// 2. CRITICAL: We must close the store when the test finishes so the
	// background sweeper goroutine doesn't leak into other tests.
	defer s.Close()

	var wg sync.WaitGroup

	// 3. Spawn 100 goroutines writing simultaneously
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			key := fmt.Sprintf("key_%d", val%10) // Overlapping keys to force collisions
			s.Set(key, []byte("data"), 5*time.Minute)
		}(i)
	}

	// 4. Spawn 100 goroutines reading simultaneously
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			key := fmt.Sprintf("key_%d", val%10)
			_ = s.Get(key)
		}(i)
	}

	// 5. Wait for all 200 goroutines to finish crashing into each other
	wg.Wait()
}
