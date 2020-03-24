package network

import (
	"github.com/hilonfot/network/conn"
	"github.com/hilonfot/network/server"
)

type Engine struct {
	server *server.Server
}

func NewEngine() *Engine {
	s := server.NewServer()
	// runtime.GOMAXPROCS(runtime.NumCPU())
	return &Engine{server: s}
}

func (e *Engine) AddRouter(msgId uint32, router conn.Router) {
	e.server.AddRouter(msgId, router)
}

func (e *Engine) Run() {
	e.server.Server()
}
