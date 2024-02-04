package server

import (
	"bufio"
	"errors"
	"log"
	"net"
	"strings"

	"github.com/Devansh3712/tandb/store"
)

const (
	// Generic commands
	CMD_GET     = "GET"
	CMD_DEL     = "DEL"
	CMD_SET     = "SET"
	CMD_KEYS    = "KEYS"
	CMD_MGET    = "MGET"
	CMD_SETEX   = "SETEX"
	CMD_EXISTS  = "EXISTS"
	CMD_EXPIRE  = "EXPIRE"
	CMD_PERSIST = "PERSIST"
	// Set commands
	CMD_SADD       = "SADD"
	CMD_SCARD      = "SCARD"
	CMD_SDIFF      = "SDIFF"
	CMD_SMEMBERS   = "SMEMBERS"
	CMD_SISMEMBER  = "SISMEMBER"
	CMD_SDIFFSTORE = "SDIFFSTORE"
)

var ErrInvalidCmd = errors.New("invalid command")

type Command struct {
	Value string
	Args  []string
	Conn  net.Conn
}

type Server struct {
	Addr     string
	Listener net.Listener
	DB       store.Store
	Commands chan Command
}

func NewServer(addr string) Server {
	return Server{
		Addr: addr, Commands: make(chan Command), DB: store.NewStore(),
	}
}

func (s *Server) Start() {
	listener, err := net.Listen("tcp", s.Addr)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	s.Listener = listener

	go s.DB.CheckTTL()
	s.HandleConnections()
}

func (s *Server) HandleConnections() {
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
	for cmd := range s.Commands {
		switch cmd.Value {
		case CMD_GET:
			s.get(cmd)
		case CMD_DEL:
			s.del(cmd)
		case CMD_SET:
			s.set(cmd)
		case CMD_KEYS:
			s.keys(cmd)
		case CMD_MGET:
			s.mGet(cmd)
		case CMD_SETEX:
			s.setEx(cmd)
		case CMD_EXISTS:
			s.exists(cmd)
		case CMD_EXPIRE:
			s.expire(cmd)
		case CMD_PERSIST:
			s.persist(cmd)
		case CMD_SADD:
			s.sAdd(cmd)
		case CMD_SCARD:
			s.sCard(cmd)
		case CMD_SDIFF:
			s.sDiff(cmd)
		case CMD_SMEMBERS:
			s.sMembers(cmd)
		case CMD_SISMEMBER:
			s.sIsMember(cmd)
		case CMD_SDIFFSTORE:
			s.sDiffStore(cmd)
		default:
			cmd.error(ErrInvalidCmd)
		}
	}
}
