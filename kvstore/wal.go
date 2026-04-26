package kvstore

import "time"

type ActionType string

const (
	ActionSet    ActionType = "SET"
	ActionDelete ActionType = "DELETE"
)

type LogEntry struct {
	Action     ActionType `json:"action"`
	Key        string     `json:"key"`
	Value      []byte     `json:"value,omitempty"`
	Expiration time.Time  `json:"expiration,omitempty"`
}
