package kvstore

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

func (s *KVStore) String() string {
	return fmt.Sprintf("KVStore: %v", s.store)
}

func NewKVStore(sweepInterval time.Duration) *KVStore {
	s := KVStore{
		store: make(map[string]Entry),
		done:  make(chan struct{}),
	}

	go s.StartSweeper(sweepInterval)

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
