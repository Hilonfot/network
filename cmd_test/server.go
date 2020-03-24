package main

import (
	"github.com/hilonfot/network"
	"github.com/hilonfot/network/conn"
	"github.com/hilonfot/network/utils/log"
)

//  CGO_ENABLED=0 GOOS=linux GOARCH=amd64  go build

type PingRouter struct {
	conn.Router
}

func (p *PingRouter) Handle(r *conn.Request) {
	log.Info("Call PingRouter Handle")
	log.Info("recv from client: msgId= ", r.GetMsgID(), ",data=", string(r.GetData()))

	err := r.GetConnection().SendBuffMsg(0, []byte("Ping Ping Ping Ping "))
	if err != nil {
		log.Error(err.Error())
	}
}

type HelloRouter struct {
	conn.Router
}

func (h *HelloRouter) Handle(r *conn.Request) {
	log.Info("Call PingRouter Handle")
	log.Info("recv from client: msgId= ", r.GetMsgID(), ",data=", string(r.GetData()))
	err := r.GetConnection().SendBuffMsg(1, []byte("Hello Router "))
	if err != nil {
		log.Error(err.Error())
	}
}

func main() {
	// 创建一个server
	s := network.NewEngine()
	// 注册回调hook函数

	// 配置路由
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})
	// start
	s.Run()
}
