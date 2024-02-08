package server

import (
	"fmt"
	"strconv"
	"time"
)

func (s *Server) get(cmd Command) {
	if len(cmd.Args) < 1 {
		cmd.error(ErrNotEnoughArgs)
		return
	}
	result, err := s.DB.Get(cmd.Args[0])
	if err != nil {
		cmd.error(err)
		return
	}
	cmd.write(string(result))
}

func (s *Server) set(cmd Command) {
	if len(cmd.Args) < 2 {
		cmd.error(ErrNotEnoughArgs)
		return
	}
	err := s.DB.Set(cmd.Args[0], []byte(cmd.Args[1]))
	if err != nil {
		cmd.error(err)
	}
}

func (s *Server) setEx(cmd Command) {
	if len(cmd.Args) < 3 {
		cmd.error(ErrNotEnoughArgs)
		return
	}
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
	if len(cmd.Args) < 1 {
		cmd.error(ErrNotEnoughArgs)
		return
	}
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
	if len(cmd.Args) < 2 {
		cmd.error(ErrNotEnoughArgs)
		return
	}
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
	if len(cmd.Args) < 1 {
		cmd.error(ErrNotEnoughArgs)
		return
	}
	ok := s.DB.Exists(cmd.Args[0])
	if !ok {
		cmd.write("FALSE")
		return
	}
	cmd.write("TRUE")
}

func (s *Server) persist(cmd Command) {
	if len(cmd.Args) < 1 {
		cmd.error(ErrNotEnoughArgs)
		return
	}
	err := s.DB.Persist(cmd.Args[0])
	if err != nil {
		cmd.error(err)
	}
}

func (s *Server) expireTime(cmd Command) {
	if len(cmd.Args) < 1 {
		cmd.error(ErrNotEnoughArgs)
		return
	}
	exp, err := s.DB.ExpireTime(cmd.Args[0])
	if err != nil {
		cmd.error(err)
		return
	}
	cmd.write(strconv.FormatInt(exp, 10))
}

func (s *Server) ttl(cmd Command) {
	if len(cmd.Args) < 1 {
		cmd.error(ErrNotEnoughArgs)
		return
	}
	ttl, err := s.DB.TTL(cmd.Args[0])
	if err != nil {
		cmd.error(err)
		return
	}
	cmd.write(strconv.FormatFloat(ttl, 'f', 0, 64))
}
