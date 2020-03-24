package globalobj

import (
	"encoding/json"
	"github.com/hilonfot/network/utils/log"
	"io/ioutil"
)

type GlobalObj struct {
	Host    string
	TcpPort int
	Name    string
	Version string // 当前服务器版本号

	MaxPacketSize uint32
	MaxConn       int

	WorkerPoolSize   uint32
	MaxWorkerTaskLen uint32 // 负责的任务队列最大任务存储数量

	ConfFilePath  string // config file path
	MaxMsgChanLen uint32 // 消息通道最大len
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	log.Info("Load GlobalObject !")
	data, err := ioutil.ReadFile(g.ConfFilePath)
	if err != nil {
		// 跳过加载文件
		log.Info(err.Error(), "，Use default Settings ！")
		return
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err.Error())
	}
}

func init() {
	GlobalObject = &GlobalObj{
		Host:             "0.0.0.0",
		TcpPort:          8888,
		Name:             "ServerApp",
		Version:          "0.0.1",
		MaxPacketSize:    4096,
		MaxConn:          12000,
		WorkerPoolSize:   8,
		MaxWorkerTaskLen: 1024,
		ConfFilePath:     "./setting.json",
		// MaxMsgChanLen:    0,
	}
	GlobalObject.Reload()
}
