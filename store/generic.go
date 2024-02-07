package store

import (
	"errors"
	"time"
)

var (
	ErrKeyExists    = errors.New("the key already exists")
	ErrKeyNotExists = errors.New("the key does not exist")
)

type Value struct {
	Timestamp  time.Time
	Data       []byte
	Expiration time.Duration
}

func (v *Value) expired() bool {
	if v.Expiration == -1 {
		return false
	}
	now := time.Now()
	ttl := v.Timestamp.Add(v.Expiration)
	return now.After(ttl)
}

// Check if a key exists.
func (s *Store) Exists(key string) bool {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	_, ok := s.Records[key]
	return ok
}

// Store a key-value pair with an expiration time (in seconds).
// If the value of expiration is -1, persist the record.
//
// Persistence refers to a key-value pair with no expiration.
func (s *Store) SetEx(key string, value []byte, expiration time.Duration) error {
	if s.Exists(key) {
		return ErrKeyExists
	}
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	s.Records[key] = Value{
		Timestamp: time.Now(), Data: value, Expiration: expiration * time.Second,
	}
	return nil
}

// Store a persistent key-value pair.
func (s *Store) Set(key string, value []byte) error {
	return s.SetEx(key, value, -1)
}

// Fetch a value of the input key.
func (s *Store) Get(key string) ([]byte, error) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	value, ok := s.Records[key]
	if !ok {
		return nil, ErrKeyNotExists
	}
	return value.Data, nil
}

// Fetch values of multiple keys at once.
// If key does not exist, <nil> is appended as the value.
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

// Delete a key-value pair.
func (s *Store) Del(key string) error {
	if !s.Exists(key) {
		return ErrKeyNotExists
	}

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	delete(s.Records, key)
	return nil
}

// Set or update the expiration time for a key-value pair.
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

// Set or update a key-value pair to persist in store.
func (s *Store) Persist(key string) error {
	return s.Expire(key, -1)
}

// Fetch all keys in the store.
func (s *Store) Keys() []string {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	var keys []string
	for key := range s.Records {
		keys = append(keys, key)
	}
	return keys
}

// Returns the expiration time of a key as Unix timestamp.
func (s *Store) ExpireTime(key string) (int64, error) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	value, ok := s.Records[key]
	if !ok {
		return 0, ErrKeyNotExists
	}
	expireTime := value.Timestamp.Add(value.Expiration)
	return expireTime.Unix(), nil
}

// Returns the expiration time in seconds of a key
func (s *Store) TTL(key string) (float64, error) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	value, ok := s.Records[key]
	if !ok {
		return 0, ErrKeyNotExists
	}
	ttl := time.Until(value.Timestamp.Add(value.Expiration))
	return ttl.Seconds(), nil
}
