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
		cmd.error(err)
		return
	}
	cmd.write(string(result))
}

func (s *Server) set(cmd Command) {
	err := s.DB.Set(cmd.Args[0], []byte(cmd.Args[1]))
	if err != nil {
		cmd.error(err)
	}
}

func (s *Server) setEx(cmd Command) {
	ttl, err := strconv.Atoi(cmd.Args[2])
	if err != nil {
		cmd.error(err)
	}
	err = s.DB.SetEx(cmd.Args[0], []byte(cmd.Args[1]), time.Duration(ttl))
	if err != nil {
		cmd.error(err)
	}
}

func (s *Server) del(cmd Command) {
	err := s.DB.Del(cmd.Args[0])
	if err != nil {
		cmd.error(err)
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
		cmd.error(err)
	}
	err = s.DB.Expire(cmd.Args[0], time.Duration(expiration))
	if err != nil {
		cmd.error(err)
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
		cmd.error(err)
	}
}

func (s *Server) sAdd(cmd Command) {
	s.DB.SAdd(cmd.Args[0], cmd.Args[1])
}

func (s *Server) sMembers(cmd Command) {
	elements, err := s.DB.SMembers(cmd.Args[0])
	if err != nil {
		cmd.error(err)
		return
	}
	for index, element := range elements {
		cmd.write(fmt.Sprintf("%d) %s", index+1, element))
	}
}

func (s *Server) sCard(cmd Command) {
	size, err := s.DB.SCard(cmd.Args[0])
	if err != nil {
		cmd.error(err)
		return
	}
	cmd.write(strconv.Itoa(size))
}

func (s *Server) sIsMember(cmd Command) {
	ok, err := s.DB.SIsMember(cmd.Args[0], cmd.Args[1])
	if err != nil {
		cmd.error(err)
		return
	}
	if !ok {
		cmd.write("FALSE")
		return
	}
	cmd.write("TRUE")
}

func (s *Server) sDiff(cmd Command) {
	elements, err := s.DB.SDiff(cmd.Args[0], cmd.Args[1])
	if err != nil {
		cmd.error(err)
		return
	}
	for index, element := range elements {
		cmd.write(fmt.Sprintf("%d) %s", index+1, element))
	}
}

func (s *Server) sDiffStore(cmd Command) {
	err := s.DB.SDiffStore(cmd.Args[0], cmd.Args[1], cmd.Args[2])
	if err != nil {
		cmd.error(err)
	}
}
