// Unsorted set implementation using map. It stores the elements
// as a key and an empty struct as value, as an empty struct takes
// 0 bytes of memory. The methods are concurrency safe, using a
// read-write mutex.

package main

import (
	"errors"
	"sync"
)

var (
	ErrElementNotExists = errors.New("the element does not exist in set")
)

func NewSet() Set {
	return Set{
		Mutex:    &sync.RWMutex{},
		Elements: make(map[string]struct{}),
	}
}

func (s *Set) Size() int {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	return len(s.Elements)
}

func (s *Set) Add(element string) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	s.Elements[element] = struct{}{}
}

func (s *Set) Exists(element string) bool {
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()

	_, ok := s.Elements[element]
	return ok
}

func (s *Set) Remove(element string) error {
	if ok := s.Exists(element); !ok {
		return ErrElementNotExists
	}

	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	delete(s.Elements, element)

	return nil
}

func (s1 *Set) Union(s2 Set) Set {
	elements := NewSet()

	s1.Mutex.RLock()
	for element := range s1.Elements {
		elements.Add(element)
	}
	s1.Mutex.RUnlock()

	s2.Mutex.RLock()
	for element := range s2.Elements {
		elements.Add(element)
	}
	s2.Mutex.RUnlock()
	return elements
}

func (s1 *Set) Intersection(s2 Set) Set {
	s1.Mutex.RLock()
	s2.Mutex.RLock()
	defer s1.Mutex.RUnlock()
	defer s2.Mutex.RUnlock()

	elements := NewSet()
	for element := range s1.Elements {
		if s2.Exists(element) {
			elements.Add(element)
		}
	}
	return elements
}

func (s1 *Set) Difference(s2 Set) Set {
	s1.Mutex.RLock()
	s2.Mutex.RLock()
	defer s1.Mutex.RUnlock()
	defer s2.Mutex.RUnlock()

	elements := NewSet()
	for element := range s1.Elements {
		if !s2.Exists(element) {
			elements.Add(element)
		}
	}
	return elements
}

func (s1 *Set) Subset(s2 Set) bool {
	s1.Mutex.RLock()
	s2.Mutex.RLock()
	defer s1.Mutex.RUnlock()
	defer s2.Mutex.RUnlock()

	for element := range s1.Elements {
		if !s2.Exists(element) {
			return false
		}
	}
	return true
}
