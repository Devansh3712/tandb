package main

import (
	"errors"
	"sync"
	"time"
)

var (
	ErrKeyExists    = errors.New("the key already exists")
	ErrKeyNotExists = errors.New("the key does not exist")
)

func NewStore() Store {
	return Store{
		Mutex:   &sync.RWMutex{},
		Records: make(map[string]Value),
	}
}

func (s *Store) Has(key string) bool {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	_, ok := s.Records[key]
	return ok
}

func (s *Store) SetEx(key string, value []byte, expiration time.Duration) error {
	if s.Has(key) {
		return ErrKeyExists
	}
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	s.Records[key] = Value{
		Timestamp: time.Now(), Data: value, Expiration: expiration,
	}
	return nil
}

func (s *Store) Set(key string, value []byte) error {
	return s.SetEx(key, value, -1)
}

func (s *Store) Get(key string) (Value, error) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	value, ok := s.Records[key]
	if !ok {
		return value, ErrKeyNotExists
	}
	return value, nil
}

// Run a background job to check if any key has reached its expiration
// time and remove it from the store
func (s *Store) checkTTL() {
	for {
		time.Sleep(time.Second)

		s.Mutex.Lock()
		for key, value := range s.Records {
			if value.expired() {
				delete(s.Records, key)
			}
		}
		s.Mutex.Unlock()
	}
}
