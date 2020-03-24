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

// 创建连接的时候执行
func DoConnBegin(conn *conn.Connection) {
	log.Info("DoConnection Begin is Called...")

	// 设置两个连接属性，在创建连接之后
	log.Info("Set conn Name, Home done")

	conn.SetProperty("Name", "hilonfot")
	conn.SetProperty("Home", "http://hilonfot.com")

	err := conn.SendMsg(2, []byte("DoConnection Begin... "))
	if err != nil {
		log.Error(err.Error())
	}
}

// 断开连接的时候执行
func DoConnLost(conn *conn.Connection) {
	if name, err := conn.GetProperty("Name"); err == nil {
		log.Info("Conn Property Name = ", name)
	}

	if home, err := conn.GetProperty("Home"); err == nil {
		log.Info("Conn Property Home = ", home)
	}

	log.Info("DoConnectionLost is Called...")

}

func main() {
	// 创建一个server
	s := network.NewEngine()
	// 注册回调hook函数
	s.SetOnConnStart(DoConnBegin)
	s.SetOnConnStop(DoConnLost)
	// 配置路由
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})
	// start
	s.Run()
}
