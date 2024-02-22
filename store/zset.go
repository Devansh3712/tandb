package store

import "github.com/Devansh3712/tandb/zset"

func (s *Store) ZAdd(set, key string) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	value, ok := s.ZSets[set]
	if !ok {
		value = zset.NewZSet()
	}
	value.Add(key)
	s.ZSets[set] = value
}

func (s *Store) ZMembers(set string) ([]string, error) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	value, ok := s.ZSets[set]
	if !ok {
		return nil, ErrSetNotExists
	}
	return value.Members(), nil
}
