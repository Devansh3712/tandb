package main

import "github.com/Devansh3712/tandb/server"

func main() {
	server := server.NewServer(":8000")
	server.Start()
}
