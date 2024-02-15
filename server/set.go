package server

import (
	"fmt"
	"strconv"
)

func (s *Server) sAdd(cmd Command) {
	if len(cmd.Args) < 2 {
		cmd.error(ErrNotEnoughArgs)
		return
	}
	s.DB.SAdd(cmd.Args[0], cmd.Args[1])
}

func (s *Server) sMembers(cmd Command) {
	if len(cmd.Args) < 1 {
		cmd.error(ErrNotEnoughArgs)
		return
	}
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
	if len(cmd.Args) < 1 {
		cmd.error(ErrNotEnoughArgs)
		return
	}
	size, err := s.DB.SCard(cmd.Args[0])
	if err != nil {
		cmd.error(err)
		return
	}
	cmd.write(strconv.Itoa(size))
}

func (s *Server) sIsMember(cmd Command) {
	if len(cmd.Args) < 2 {
		cmd.error(ErrNotEnoughArgs)
		return
	}
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
	if len(cmd.Args) < 2 {
		cmd.error(ErrNotEnoughArgs)
		return
	}
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
	if len(cmd.Args) < 3 {
		cmd.error(ErrNotEnoughArgs)
		return
	}
	err := s.DB.SDiffStore(cmd.Args[0], cmd.Args[1], cmd.Args[2])
	if err != nil {
		cmd.error(err)
	}
}

func (s *Server) sInter(cmd Command) {
	if len(cmd.Args) < 2 {
		cmd.error(ErrNotEnoughArgs)
		return
	}
	elements, err := s.DB.SInter(cmd.Args[0], cmd.Args[1])
	if err != nil {
		cmd.error(err)
		return
	}
	for index, element := range elements {
		cmd.write(fmt.Sprintf("%d) %s", index+1, element))
	}
}
