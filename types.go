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
	if v.Expiration == -1 {
		return false
	}
	now := time.Now()
	TTL := v.Timestamp.Add(v.Expiration)
	return now.After(TTL)
}
