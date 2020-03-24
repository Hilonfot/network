package server

import (
	"fmt"
	"github.com/hilonfot/network/conn"
	"github.com/hilonfot/network/server/globalobj"
	"github.com/hilonfot/network/utils/log"
	"net"
	"time"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int

	// 当前server的消息管理模块，用来绑定MsgId和对应的处理方法
	MsgHandle *conn.MsgHandle

	// 当前服务的连接管理器
	ConnMgr *conn.ConnManager
}

func NewServer() *Server {
	return &Server{
		Name:      globalobj.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        globalobj.GlobalObject.Host,
		Port:      globalobj.GlobalObject.TcpPort,
		MsgHandle: conn.NewMsgHandle(),
		ConnMgr:   conn.NewConnManager(),
	}
}

func (s *Server) Server() {
	s.Start()

	// 阻塞，否则主程序退出，ListenTCP的goroutine将退出
	for {
		time.Sleep(10 * time.Second)
	}
}

func (s *Server) Start() {
	log.Infof("[START] Server name: %s listenner at IP : %s, Port %d, is starting \n", s.Name, s.IP, s.Port)
	log.Infof("[NetWork] Version: %s, MaxConn: %d, MaxPacketSize: %d \n",
		globalobj.GlobalObject.Version,
		globalobj.GlobalObject.MaxConn,
		globalobj.GlobalObject.MaxPacketSize,
	)
	// 开始一个goroutine去做服务端的监听业务
	go func() {
		// 开启worker工作池机制
		s.MsgHandle.StartWorkerPool()

		// 1.获取一个TCP的addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))

		if err != nil {
			log.Error("resolve tcp addr err: ", err.Error())
			return
		}

		// 2.监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			log.Error("listen", s.IPVersion, "err", err.Error())
			panic("程序监听错误 >>>>>  ")
		}
		// 已经监听成功
		log.Info("start net server ", s.Name, " success, now listening...")

		var cid uint32

		// 3.启动server网络连接业务
		for {
			// 3.1 阻塞等待客户端建立连接请求
			tcpConn, err := listener.AcceptTCP()
			if err != nil {
				log.Error("Accept err ", err.Error())
				continue
			}

			// TODO Server.Start() 设置服务最大连接控制，如果超过最大连接，那么关闭此新的连接
			if s.ConnMgr.Len() >= globalobj.GlobalObject.MaxConn {
				tcpConn.Close()
				continue
			}

			// 处理新连接请求的业务方法，此时应该有handler和conn是绑定的
			dealConn := conn.NewConnection(s.ConnMgr, tcpConn, cid, s.MsgHandle)

			// 分层设计，只能从上到下操作，不能互相调用
			s.ConnMgr.Add(dealConn)

			cid++
			// 启动当前连接的业务处理,处理数据读写分离
			go dealConn.Start()

		}
	}()
}

// 返回服务连接管理
func (s *Server) GetConnMgr() *conn.ConnManager {
	return s.ConnMgr
}

// 停止Server,将连接管理器全部Clear
func (s *Server) Stop() {
	log.Info("[STOP] Zinx server, name ", s.Name)
	s.ConnMgr.ClearConn()
	// TODO goroutine 关闭处理，关闭主协程的等待
}

// 路由功能： 给当前服务注册一个路由器服务方法，供客户端处理使用
func (s *Server) AddRouter(msgId uint32, router conn.Router) {
	s.MsgHandle.AddRouter(msgId, router)
	log.Info("Add Router success!")
}

// 设置hook函数 >>>>>>
