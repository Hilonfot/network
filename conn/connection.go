package conn

import (
	"errors"
	"github.com/hilonfot/network/message"
	"github.com/hilonfot/network/server/globalobj"
	"github.com/hilonfot/network/utils/catch"
	"github.com/hilonfot/network/utils/log"
	"io"
	"net"
	"sync"
)

type Connection struct {
	// ConnMgr
	ConnMgr *ConnManager
	// 当前连接的socket tcp套节字
	Conn *net.TCPConn
	// 当前连接的ID 也可称为SessionID,ID全局唯一
	ConnID uint32
	// 当前连接的状态
	IsClosed bool
	// 该连接的处理方式router
	MsgHanler *MsgHandle

	// 告知该连接已退出 / 停止的channel
	ExitBuffChan chan bool
	// 给缓冲队列发送数据的channel
	// 如果向缓冲队列发送数据，那么把数据发送到这个channel下
	// SendBuffChan chan []byte

	// 无缓冲通道，用于读 写 两个goroutine之间的消息通信
	MsgChan chan []byte
	// 有缓冲通道
	MsgBuffChan chan []byte

	// 连接属性
	Property map[string]interface{}

	// 保护连接属性的读写锁
	mx sync.RWMutex
}

// 创建连接
func NewConnection(ConnMgr *ConnManager, conn *net.TCPConn, connID uint32, msgHandler *MsgHandle) *Connection {
	c := &Connection{
		ConnMgr:      ConnMgr,
		Conn:         conn,
		ConnID:       connID,
		IsClosed:     false,
		MsgHanler:    msgHandler,
		ExitBuffChan: make(chan bool, 1),
		MsgChan:      make(chan []byte),
		MsgBuffChan:  make(chan []byte, globalobj.GlobalObject.MaxMsgChanLen),
		Property:     make(map[string]interface{}),
	}
	return c
}

// 启动连接，让当前连接开始工作
func (c *Connection) Start() {

	// 开启处理该连接取到客户端数据之后的请求业务
	go c.StartReader()
	go c.StartWrite()

	c.ConnMgr.CallOnConnStart(c)

	for {
		select {
		case <-c.ExitBuffChan:
			// 得到退出消息，不再阻塞
			return
		}
	}
}

// 读写分离
func (c *Connection) StartWrite() {
	// 捕获panic
	defer catch.PanicHandler()

	log.Info("Writer Goroutines  is running")
	defer log.Info(c.RemoteAddr().String(), " conn writer exit!")

	for {
		select {
		case data := <-c.MsgChan:
			if _, err := c.Conn.Write(data); err != nil {
				log.Info("Send data error: ", err, " Conn writer exit")
				return
			}
		case data, ok := <-c.MsgBuffChan:
			if ok {
				if _, err := c.Conn.Write(data); err != nil {
					log.Info("Send data error: ", err, " Conn writer exit")
					return
				}
			} else {
				break
				log.Error("msgBuffChan is Closed")
			}
		case <-c.ExitBuffChan:
			return
		}
	}
}

func (c *Connection) StartReader() {
	// 捕获panic
	defer catch.PanicHandler()

	log.Info("Reader Goroutine is running")
	defer log.Info(c.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()
	for {
		// 创建拆包解包对象
		dp := message.NewDataPack()

		// 读取客户端的msg head
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnetcion(), headData); err != nil {
			log.Info("read msg head error ", err.Error())
			c.ExitBuffChan <- true
			break
		}

		// 拆包，得到msgId和dataLen放在msg中
		msg, err := dp.UpPack(headData)
		if err != nil {
			log.Info("unpack error ", err)
			c.ExitBuffChan <- true
			break
		}

		// 根据dataLen读取data，放到msg.Data中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnetcion(), data); err != nil {
				c.ExitBuffChan <- true
				break
			}
		}
		msg.SetData(data)

		// 根据当前客户端请求的request数据
		req := &Request{
			conn: c,
			msg:  msg,
		}

		if globalobj.GlobalObject.WorkerPoolSize > 0 {
			c.MsgHanler.SendMsgToTaskQueue(req)
		} else {
			// 从路由Routers中找到注册绑定conn的对应方法
			go c.MsgHanler.DoMsgHandler(req)
		}
	}
}

// 获取远程客户端的地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 停止连接，结束当前连接状态
func (c *Connection) Stop() {
	log.Info("Conn stop ... ConnID=", c.ConnID)

	// 如果当前连接已经关闭
	if c.IsClosed == true {
		return
	}

	c.IsClosed = true

	// 调用Stop  Hook 回调函数
	c.ConnMgr.CallOnConnStop(c)

	// 关闭socket连接
	c.Conn.Close()

	// 通知从缓冲队列读取数据的业务，该连接已经关闭
	c.ExitBuffChan <- true

	// 将连接从连接管理器中删除
	c.ConnMgr.Remove(c)

	// 关闭该连接的全部管道
	close(c.ExitBuffChan)
	close(c.MsgBuffChan)
	close(c.MsgChan)
}

// 从当前连接获取原始的socket TCPconn
func (c *Connection) GetTCPConnetcion() *net.TCPConn {
	return c.Conn
}

// 获取当前连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// 设置连接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.Property[key] = value
}

// 获取连接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	if value, ok := c.Property[key]; ok {
		return value, nil
	}
	return nil, errors.New("no propery found")
}

// 移除连接属性
func (c *Connection) RemoveProperty(key string) {
	c.mx.Lock()
	defer c.mx.Unlock()

	delete(c.Property, key)
}

// >>>>>>消息处理 >>>>>
// 直接将数据发送给远程的TCP客户端
func (c *Connection) Send(data []byte) error {
	// TODO  saomethings...
	return nil
}

// 将数据发送给缓冲队列，通过专门从缓冲队列读取数据的goroutine写给客户端
func (c *Connection) SendBuff() error {
	// TODO  saomethings...
	return nil
}

// 发送不带缓冲的消息给客户端
func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.IsClosed == true {
		return errors.New("Connection closed when send msg ")
	}

	// 将data封包
	dp := message.NewDataPack()
	msg, err := dp.Pack(message.NewMessage(msgId, data))
	if err != nil {
		log.Error("Pack error msg id = ", msgId)
		return errors.New("pack error msg ")
	}

	// 写入客户端
	c.MsgChan <- msg // 将直接回写给conn.write的方法，改为发送给channel，供writer读取

	return nil
}

// 发送带缓冲的消息给客户端
func (c *Connection) SendBuffMsg(msgId uint32, data []byte) error {
	if c.IsClosed == true {
		return errors.New("Connection closed when send msg ")
	}

	// 将data封包
	dp := message.NewDataPack()
	msg, err := dp.Pack(message.NewMessage(msgId, data))
	if err != nil {
		log.Error("Pack error msg id = ", msgId)
		return errors.New("pack error msg ")
	}

	// 写入客户端
	c.MsgBuffChan <- msg // 将直接回写给conn.write的方法，改为发送给channel，供writer读取

	return nil
}
