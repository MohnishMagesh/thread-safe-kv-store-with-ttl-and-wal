package kvstore

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

func (s *KVStore) Recovery(file *os.File) error {
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var entry LogEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			return fmt.Errorf("failed to unmarshal log entry: %w", err)
		}

		switch entry.Action {
		case ActionSet:
			s.store[entry.Key] = Entry{
				bytes:          entry.Value,
				expirationTime: entry.Expiration,
			}
		case ActionDelete:
			delete(s.store, entry.Key)
		default:
			return fmt.Errorf("unknown action type: %s", entry.Action)
		}
	}

	return scanner.Err()
}

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

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.writer != nil {
		s.writer.Flush()
	}
	if s.file != nil {
		s.file.Close()
	}
}
