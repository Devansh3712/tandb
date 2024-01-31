package main

import (
	"net"
	"sync"
	"time"
)

type Value struct {
	Timestamp  time.Time
	Data       []byte
	Expiration time.Duration
}

type Set struct {
	Mutex    *sync.RWMutex
	Elements map[string]struct{}
}

type Node struct {
	Color  int
	Value  string
	Score  int
	Parent *Node
	Left   *Node
	Right  *Node
}

type RBTree struct {
	Mutex *sync.RWMutex
	Root *Node
}

type Store struct {
	Mutex   *sync.RWMutex
	Records map[string]Value
	Sets    map[string]Set
}

type Command struct {
	Value string
	Args  []string
	Conn  net.Conn
}

type Server struct {
	Wg       sync.WaitGroup
	Addr     string
	Listener net.Listener
	Commands chan Command
	DB       Store
}

func (v *Value) expired() bool {
	if v.Expiration == -1 {
		return false
	}
	now := time.Now()
	ttl := v.Timestamp.Add(time.Second * v.Expiration)
	return now.After(ttl)
}
