package kvstore

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
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

	file    *os.File
	writer  *bufio.Writer
	encoder *json.Encoder
}

func (s *KVStore) String() string {
	return fmt.Sprintf("KVStore: %v", s.store)
}

func NewKVStore(sweepInterval time.Duration) (*KVStore, error) {
	s := KVStore{
		store: make(map[string]Entry),
		done:  make(chan struct{}),
	}

	file, err := os.OpenFile("wal.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open wal: %w", err)
	}
	s.file = file
	s.writer = bufio.NewWriter(s.file)
	s.encoder = json.NewEncoder(s.writer)

	if err := s.Recovery(file); err != nil {
		return nil, fmt.Errorf("failed to recover from wal: %w", err)
	}

	go s.StartSweeper(sweepInterval)

	return &s, err
}

func (s *KVStore) Get(key string) Entry {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.store[key]
}

func (s *KVStore) Set(key string, val []byte, ttl time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	logEntry := LogEntry{
		Action:     ActionSet,
		Key:        key,
		Value:      val,
		Expiration: time.Now().Add(ttl),
	}
	s.encoder.Encode(logEntry)

	s.store[key] = Entry{bytes: val, expirationTime: time.Now().Add(ttl)}
}

func (s *KVStore) Delete(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	LogEntry := LogEntry{
		Action: ActionDelete,
		Key:    key,
	}
	s.encoder.Encode(LogEntry)

	delete(s.store, key)
}
