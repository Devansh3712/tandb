package zset

import (
	"errors"
	"sync"
)

var ErrElementNotExists = errors.New("the element does not exist in set")

type ZSet struct {
	Mutex    *sync.RWMutex
	Elements *RBTree
}

func NewZSet() ZSet {
	return ZSet{
		Mutex:    &sync.RWMutex{},
		Elements: NewRBTree(),
	}
}

func (z *ZSet) Size() int {
	z.Mutex.RLock()
	defer z.Mutex.RUnlock()

	return z.Elements.Count
}

func (z *ZSet) Add(element string) {
	z.Mutex.Lock()
	defer z.Mutex.Unlock()

	z.Elements.insert(element)
}

func (z *ZSet) Exists(element string) bool {
	z.Mutex.RLock()
	defer z.Mutex.RUnlock()

	_, ok := z.Elements.search(element)
	return ok
}

func (z *ZSet) Remove(element string) error {
	if ok := z.Exists(element); !ok {
		return ErrElementNotExists
	}

	z.Mutex.Lock()
	defer z.Mutex.Unlock()
	z.Elements.delete(element)

	return nil
}

func (z *ZSet) Members() []string {
	z.Mutex.RLock()
	defer z.Mutex.RUnlock()

	return z.Elements.members()
}
