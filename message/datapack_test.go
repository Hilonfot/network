package message

import (
	"io"
	"net"
	"testing"
	"time"
)

func TestDataPack(t *testing.T) {
	// 模拟服务器
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		t.Fatal(err.Error())
	}

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				t.Fatal(err.Error())
				return
			}

			// 开goroutine处理conn
			go func(conn net.Conn) {
				// 只是为了获取结构体方法，DataPack struct 为空
				dp := NewDataPack()
				// 一个conn 发送多个包，就需要循环读取包结构
				for {
					// 获取头部数据
					headData := make([]byte, dp.GetHeadLen())
					if _, err := io.ReadFull(conn, headData); err != nil {
						t.Fatal(err.Error())
						break
					}

					// 将headData upPack
					msg, err := dp.UpPack(headData)
					if err != nil {
						t.Fatal(err.Error())
						return
					}

					// 读取data内容
					if msg.DataLen > 0 {
						data := make([]byte, msg.DataLen)
						if _, err := io.ReadFull(conn, data); err != nil {
							t.Fatal(err.Error())
							return
						}
						msg.SetData(data)
						t.Logf("recv Len:%d ID:%d data:%v \n", msg.GetDataLen(), msg.GetMsgId(), string(msg.GetData()))
					}
				}
			}(conn)
		}
	}()

	time.Sleep(time.Second * 1)

	// 模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	defer conn.Close()
	if err != nil {
		t.Log(err.Error())
		return
	}
	dp := NewDataPack()
	for i := 0; i < 5; i++ {

		msgClient := NewMessage(uint32(i), []byte("你好呀"))
		t.Log(msgClient)
		binaryData, err := dp.Pack(msgClient)
		if err != nil {
			t.Fatal(err.Error())
			return
		}
		t.Log(binaryData)
		conn.Write(binaryData)
	}

	time.Sleep(time.Second * 5)
	// select {}
}
