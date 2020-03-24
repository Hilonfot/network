package main

import (
	"github.com/hilonfot/network/message"
	"github.com/hilonfot/network/utils/log"
	"io"
	"net"
	"time"
)

func main() {
	// 客户端拨通连接
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		panic("Dial err !")
	}

	// 模拟客户端封包拆包message
	for {
		// 发送封包消息
		dp := message.NewDataPack()
		msg, _ := dp.Pack(message.NewMessage(0, []byte("你好")))
		_, err := conn.Write(msg)
		if err != nil {
			log.Error("write error err ", err.Error())
			return
		}

		// 先读取出流中的head部分
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			log.Error(err.Error())
			// 一般情况是出现EOF出错，直接break出当前循环
			break
		}

		// 将headData字节流，拆包到msg中
		msgRecv, err := dp.UpPack(headData)
		if err != nil {
			log.Error(err.Error())
			return
		}
		if msgRecv.GetDataLen() > 0 {
			// 如果有data数据， 需要继续从conn里面读取
			dataRecv := make([]byte, msgRecv.GetDataLen())

			// 填充满读取
			_, err := io.ReadFull(conn, dataRecv)
			if err != nil {
				log.Error("server unpack data err:", err.Error())
				return
			}
			msgRecv.SetData(dataRecv)
		}

		log.Info("==> Recv Msg: ID=", msgRecv.Id, ", len=", msgRecv.DataLen, ", data=", string(msgRecv.Data))

		log.Info("Sleep 4 second >>>>>>> ")
		time.Sleep(4 * time.Second)
	}
}
