package main

import (
	"errors"
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
	c.write(err.Error())
}

func (s *Server) get(cmd Command) {
	result, err := s.DB.Get(cmd.Args[0])
	if err != nil {
		cmd.error(ErrKeyNotExists)
		return
	}
	cmd.write(string(result.Data))
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
