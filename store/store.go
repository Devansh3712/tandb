package store

import (
	"sync"
	"time"

	"github.com/Devansh3712/tandb/set"
)

type Store struct {
	Mutex   *sync.RWMutex
	Records map[string]Value
	Sets    map[string]set.Set
}

func NewStore() Store {
	return Store{
		Mutex:   &sync.RWMutex{},
		Records: make(map[string]Value),
		Sets:    make(map[string]set.Set),
	}
}

// Run a background job to check if any key has reached its expiration
// time and remove it from the store.
func (s *Store) CheckTTL() {
	for {
		time.Sleep(time.Second)

		for key, value := range s.Records {
			if value.expired() {
				s.Del(key)
			}
		}
	}
}
