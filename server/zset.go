package server

import "fmt"

func (s *Server) zAdd(cmd Command) {
	if len(cmd.Args) < 2 {
		cmd.error(ErrNotEnoughArgs)
		return
	}
	s.DB.ZAdd(cmd.Args[0], cmd.Args[1])
}

func (s *Server) zMembers(cmd Command) {
	if len(cmd.Args) < 1 {
		cmd.error(ErrNotEnoughArgs)
		return
	}
	elements, err := s.DB.ZMembers(cmd.Args[0])
	if err != nil {
		cmd.error(err)
		return
	}
	for index, element := range elements {
		cmd.write(fmt.Sprintf("%d) %s", index+1, element))
	}
}
