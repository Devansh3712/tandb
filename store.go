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
		Sets: make(map[string]Set),
	}
}

func (s *Store) Exists(key string) bool {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	_, ok := s.Records[key]
	return ok
}

func (s *Store) SetEx(key string, value []byte, expiration time.Duration) error {
	if s.Exists(key) {
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

func (s *Store) Get(key string) ([]byte, error) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	value, ok := s.Records[key]
	if !ok {
		return nil, ErrKeyNotExists
	}
	return value.Data, nil
}

func (s *Store) MGet(keys []string) [][]byte {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	var values [][]byte

	for _, key := range keys {
		value, _ := s.Get(key)
		values = append(values, value)
	}
	return values
}

func (s *Store) Del(key string) error {
	if !s.Exists(key) {
		return ErrKeyNotExists
	}
	delete(s.Records, key)
	return nil
}

func (s *Store) Expire(key string, expiration time.Duration) error {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	value, ok := s.Records[key]
	if !ok {
		return ErrKeyNotExists
	}
	value.Expiration = expiration
	s.Records[key] = value
	return nil
}

func (s *Store) Persist(key string) error {
	return s.Expire(key, -1)
}

func (s *Store) Keys() []string {
	var keys []string
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	for key := range s.Records {
		keys = append(keys, key)
	}
	return keys
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
