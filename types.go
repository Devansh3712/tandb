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

type Store struct {
	Mutex   *sync.RWMutex
	Records map[string]Value
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
