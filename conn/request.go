package conn

import "github.com/hilonfot/network/message"

type Request struct {
	// 已经和客户端建立好的连接
	conn *Connection
	// 客户端请求的数据
	msg *message.Message
}

// 获取请求连接信息
func (r *Request) GetConnection() *Connection {
	return r.conn
}

// 获取请求消息的数据
func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

// 获取请求的消息ID
func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}
