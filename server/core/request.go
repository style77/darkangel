package core

type Request struct {
	Server      *Server
	CommandName string
	Args        []string
}
