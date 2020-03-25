package conn

import (
	"errors"
	"github.com/hilonfot/network/utils/log"
	"sync"
)

// 连接管理模块
type ConnManager struct {
	connections map[uint32]*Connection // 管理的连接信息

	// 当前server 创建连接时的hook函数
	OnConnStart func(conn *Connection)

	// 当前server 断开连接时的hook函数
	OnConnStop func(conn *Connection)

	mx sync.RWMutex
}

// new
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]*Connection),
	}
}

// 获取目前所有的连接个数
func (c ConnManager) Len() int {
	return len(c.connections)
}

// 添加连接
func (c ConnManager) Add(conn *Connection) {
	c.mx.Lock()
	defer c.mx.Unlock()

	// 将conn连接添加到ConnManager中
	c.connections[conn.GetConnID()] = conn
	log.Infof("connection add to ConnManager successfully: conn num = %v", c.Len())
}

// 删除连接
func (c *ConnManager) Remove(conn *Connection) {
	c.mx.Lock()
	defer c.mx.Unlock()

	delete(c.connections, conn.GetConnID())
	log.Infof("connectoion remove ConnID= %v success: conn num = %v", conn.GetConnID(), c.Len())
}

// 利用connID获取连接
func (c *ConnManager) Get(connID uint32) (*Connection, error) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	if conn, ok := c.connections[connID]; ok {
		return conn, nil
	}
	return nil, errors.New("connection not found")
}

// 清除并停止所有连接
func (c *ConnManager) ClearConn() {
	c.mx.Lock()
	defer c.mx.Unlock()

	for connID, conn := range c.connections {
		conn.Stop()
		delete(c.connections, connID)
	}
	log.Info("Clear All Connections successfully : conn num = ", c.Len())
}

// 设置hook函数 >>>>>>

// 设置server 连接创建时的hook函数
func (c *ConnManager) SetOnConnStart(hookFunc func(conn *Connection)) {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.OnConnStart = hookFunc
}

// 设置server 断开连接时的hook函数
func (c *ConnManager) SetOnConnStop(hook func(conn *Connection)) {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.OnConnStop = hook
}

// 调用连接时OnConnStart Hook函数
func (c *ConnManager) CallOnConnStart(conn *Connection) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	if c.OnConnStart != nil {
		log.Info("--->Call On Conn Start...")
		c.OnConnStart(conn)
	}
}

// 调用断开时OnConnStop Hook 函数
func (c *ConnManager) CallOnConnStop(conn *Connection) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	if c.OnConnStop != nil {
		log.Info("---> CallOnConnStop...")
		c.OnConnStop(conn)
	}
}
