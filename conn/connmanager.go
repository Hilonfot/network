package conn

import (
	"errors"
	"github.com/hilonfot/network/utils/log"
	"sync"
)

// 连接管理模块
type ConnManager struct {
	connections map[uint32]*Connection // 管理的连接信息
	mx          sync.RWMutex
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
