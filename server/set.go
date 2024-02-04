package server

import (
	"fmt"
	"strconv"
)

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
