package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"
)

var (
	ErrInvalidExp = errors.New("invalid expiration time")
	ErrInvalidCmd = errors.New("invalid command")
)

func (c *Command) write(msg string) {
	_, err := c.Conn.Write([]byte(msg + "\n"))
	if err != nil {
		log.Printf("unable to write to connection: %v", err)
	}
}

func (c *Command) error(err error) {
	c.write("[ERROR] " + err.Error())
}

func (s *Server) get(cmd Command) {
	result, err := s.DB.Get(cmd.Args[0])
	if err != nil {
		cmd.error(ErrKeyNotExists)
		return
	}
	cmd.write(string(result))
}

func (s *Server) set(cmd Command) {
	err := s.DB.Set(cmd.Args[0], []byte(cmd.Args[1]))
	if err != nil {
		cmd.error(ErrKeyExists)
	}
}

func (s *Server) setEx(cmd Command) {
	ttl, err := strconv.Atoi(cmd.Args[2])
	if err != nil {
		cmd.error(ErrInvalidExp)
	}
	err = s.DB.SetEx(cmd.Args[0], []byte(cmd.Args[1]), time.Duration(ttl))
	if err != nil {
		cmd.error(ErrKeyExists)
	}
}

func (s *Server) del(cmd Command) {
	err := s.DB.Del(cmd.Args[0])
	if err != nil {
		cmd.error(ErrKeyNotExists)
	}
}

func (s *Server) mGet(cmd Command) {
	result := s.DB.MGet(cmd.Args)
	for index, value := range result {
		cmd.write(fmt.Sprintf("%d) %s", index+1, value))
	}
}

func (s *Server) expire(cmd Command) {
	expiration, err := strconv.Atoi(cmd.Args[1])
	if err != nil {
		cmd.error(ErrInvalidExp)
	}
	err = s.DB.Expire(cmd.Args[0], time.Duration(expiration))
	if err != nil {
		cmd.error(ErrKeyNotExists)
	}
}

func (s *Server) keys(cmd Command) {
	keys := s.DB.Keys()
	for index, key := range keys {
		cmd.write(fmt.Sprintf("%d) %s", index+1, key))
	}
}

func (s *Server) exists(cmd Command) {
	ok := s.DB.Exists(cmd.Args[0])
	if !ok {
		cmd.write("FALSE")
		return
	}
	cmd.write("TRUE")
}

func (s *Server) persist(cmd Command) {
	err := s.DB.Persist(cmd.Args[0])
	if err != nil {
		cmd.error(ErrKeyNotExists)
	}
}
