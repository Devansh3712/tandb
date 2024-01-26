package main

import (
	"bufio"
	"log"
	"net"
	"strings"
)

const (
	CMD_GET   = "GET"
	CMD_DEL   = "DEL"
	CMD_EXP   = "EXP"
	CMD_SET   = "SET"
	CMD_KEYS  = "KEYS"
	CMD_SETEX = "SETEX"
)

func NewServer(addr string) Server {
	return Server{
		Addr: addr, Commands: make(chan Command), DB: NewStore(),
	}
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	s.Listener = listener

	go s.DB.checkTTL()
	s.Wg.Add(1)
	go s.HandleConnections()
	s.Wg.Wait()
}

func (s *Server) HandleConnections() {
	defer s.Wg.Done()
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			log.Printf("unable to accept connection: %v", err)
			continue
		}
		go s.ReadCommand(conn)
		go s.HandleCommand()
	}
}

func (s *Server) ReadCommand(conn net.Conn) {
	defer conn.Close()
	for {
		input, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Printf("unable to read from connection: %v", err)
			break
		}
		cmd := strings.Trim(input, "\r\n")
		args := strings.Split(cmd, " ")
		s.Commands <- Command{
			Value: args[0], Args: args[1:], Conn: conn,
		}
	}
}

func (s *Server) HandleCommand() {
	for {
		cmd := <-s.Commands

		switch cmd.Value {
		case CMD_GET:
			s.get(cmd)
		case CMD_DEL:
			s.del(cmd)
		case CMD_EXP:
			s.exp(cmd)
		case CMD_SET:
			s.set(cmd)
		case CMD_KEYS:
			s.keys(cmd)
		case CMD_SETEX:
			s.setEx(cmd)
		default:
			cmd.write(ErrInvalidCmd.Error())
		}
	}
}
