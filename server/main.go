package main

import (
	"darkangel/server/core"
)

func main() {
	server := core.New(&core.Config{
		Host: "localhost",
		Port: "2306",
	})
	server.Run()
}
