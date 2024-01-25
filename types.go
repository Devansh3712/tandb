package main

import (
	"sync"
	"time"
)

type Value struct {
	Timestamp  time.Time
	Data       []byte
	Expiration time.Duration
}

type Store struct {
	Mutex   sync.RWMutex
	Records map[string]Value
}

func (v *Value) expired() bool {
	now := time.Now()
	TTL := v.Timestamp.Add(v.Expiration)
	if now.After(TTL) {
		return true
	}
	return false
}
