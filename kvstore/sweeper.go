package kvstore

import (
	"fmt"
	"time"
)

func (s *KVStore) StartSweeper(sweepInterval time.Duration) {
	ticker := time.NewTicker(sweepInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.done:
			fmt.Println("Sweeper shutting down...")
			return
		case t := <-ticker.C:
			fmt.Println("Sweeping at ", t)
			s.sweep()
		}
	}
}

func (s *KVStore) sweep() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for key, entry := range s.store {
		if time.Now().After(entry.expirationTime) {
			delete(s.store, key)
		}
	}
}

func (s *KVStore) Close() {
	close(s.done)
}
