package main

import (
	"fmt"
	"sync"
)

type KVStore struct {
	store map[string][]byte
	mutex sync.RWMutex
}

func (s *KVStore) Get(key string) []byte {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.store[key]
}

func (s *KVStore) Set(key string, val []byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.store[key] = val
}

func (s *KVStore) Delete(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.store, key)
}

func main() {
	s := KVStore{
		store: make(map[string][]byte),
	}

	s.Set("Ronaldo", []byte{'C', 'R', '7'})
	fmt.Println(string(s.Get("Ronaldo")))
}
