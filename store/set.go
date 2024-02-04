package store

import (
	"errors"
	"fmt"

	Set "github.com/Devansh3712/tandb/set"
)

var ErrSetNotExists = errors.New("the set does not exist")

func (s *Store) SAdd(set, key string) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	value, ok := s.Sets[set]
	// Initialize a new set if it does not exist
	if !ok {
		value = Set.NewSet()
	}
	value.Add(key)
	s.Sets[set] = value
}

func (s *Store) SMembers(set string) ([]string, error) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	var elements []string
	value, ok := s.Sets[set]
	if !ok {
		return nil, ErrSetNotExists
	}
	for element := range value.Elements {
		elements = append(elements, element)
	}
	return elements, nil
}

func (s *Store) SCard(set string) (int, error) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	value, ok := s.Sets[set]
	if !ok {
		return 0, ErrSetNotExists
	}
	return value.Size(), nil
}

func (s *Store) SIsMember(set, key string) (bool, error) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	value, ok := s.Sets[set]
	if !ok {
		return false, ErrSetNotExists
	}
	return value.Exists(key), nil
}

func (s *Store) SDiff(s1, s2 string) ([]string, error) {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	var elements []string
	v1, ok := s.Sets[s1]
	if !ok {
		return nil, fmt.Errorf("the set %s does not exist", s1)
	}
	v2, ok := s.Sets[s2]
	if !ok {
		return nil, fmt.Errorf("the set %s does not exist", s2)
	}

	diff := v1.Difference(v2)
	for element := range diff.Elements {
		elements = append(elements, element)
	}
	return elements, nil
}

func (s *Store) SDiffStore(s1, s2, s3 string) error {
	elements, err := s.SDiff(s1, s2)
	if err != nil {
		return err
	}
	for _, element := range elements {
		s.SAdd(s3, element)
	}
	return nil
}
