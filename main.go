package main

import (
	"fmt"
	"sync"
	"time"
)

type Entry struct {
	bytes          []byte
	expirationTime time.Time
}

func (e Entry) String() string {
	return fmt.Sprintf("Value: %v, Expiration Time: %v", e.bytes, e.expirationTime)
}

type KVStore struct {
	store map[string]Entry
	mutex sync.RWMutex
	done  chan struct{}
}

func (s *KVStore) startSweeper(sweepInterval time.Duration) {
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

func NewKVStore(sweepInterval time.Duration) *KVStore {
	s := KVStore{
		store: make(map[string]Entry),
		done:  make(chan struct{}),
	}

	go s.startSweeper(sweepInterval)

	return &s
}

func (s *KVStore) Get(key string) Entry {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.store[key]
}

func (s *KVStore) Set(key string, val []byte, ttl time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.store[key] = Entry{bytes: val, expirationTime: time.Now().Add(ttl)}
}

func (s *KVStore) Delete(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.store, key)
}

func (s *KVStore) Close() {
	close(s.done)
}

func main() {
	s := NewKVStore(10 * time.Second)
	defer s.Close()

	footballPlayers := []string{"Ronaldo", "Messi", "Neymar", "Mbappe", "Salah", "Kane", "Lewandowski", "De Bruyne", "Modric", "Van Dijk"}
	bytes := [][]byte{
		{'C', 'R', '7'},
		{'L', 'M', '1', '0'},
		{'N', 'J', '1', '0'},
		{'K', 'M', '7'},
		{'M', 'S', '1', '1'},
		{'H', 'K', '1', '0'},
		{'R', 'L', '9'},
		{'D', 'B', '1', '0'},
		{'L', 'M', '1', '0'},
		{'V', 'D', '1', '0'},
	}

	for i, player := range footballPlayers {
		s.Set(player, bytes[i], 5*time.Second)
		time.Sleep(2 * time.Second)
		fmt.Println(s.Get(player))
	}

	fmt.Println("Final store state:", s.store)
}
