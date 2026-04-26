package main

import (
	"fmt"
	"sync"
	"testing"
)

func TestKVStoreRaceCondition(t *testing.T) {
	// 1. Initialize our unprotected store
	s := KVStore{
		store: make(map[string][]byte),
	}

	var wg sync.WaitGroup

	// 2. Spawn 100 goroutines that all try to WRITE at the exact same time
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			key := fmt.Sprintf("key_%d", val%10) // Overlapping keys
			s.Set(key, []byte("data"))
		}(i)
	}

	// 3. Spawn 100 goroutines that all try to READ at the exact same time
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			key := fmt.Sprintf("key_%d", val%10)
			_ = s.Get(key)
		}(i)
	}

	// 4. Wait for all 200 goroutines to finish crashing into each other
	wg.Wait()
}
