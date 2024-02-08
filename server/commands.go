package server

import (
	"errors"
	"log"
)

var ErrNotEnoughArgs = errors.New("not enough arguments for command")

func (c *Command) write(msg string) {
	_, err := c.Conn.Write([]byte(msg + "\n"))
	if err != nil {
		log.Printf("unable to write to connection: %v", err)
	}
}

func (c *Command) error(err error) {
	c.write("[ERROR] " + err.Error())
}
