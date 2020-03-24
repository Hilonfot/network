package conn

import (
	"github.com/hilonfot/network/server/globalobj"
	"github.com/hilonfot/network/utils/log"
	"strconv"
)

type MsgHandle struct {
	// msgID 所对应的处理路由
	Apis map[uint32]Router
	// Worker池的开启数量
	WorkerPoolSize uint32
	// Worker负责任务的消息队列
	TaskQueue []chan *Request
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]Router),
		WorkerPoolSize: globalobj.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan *Request, globalobj.GlobalObject.WorkerPoolSize),
	}
}

// 为消息添加具体的处理逻辑
func (m *MsgHandle) AddRouter(msgId uint32, router Router) {
	// 判断当前msg的API处理方法是否已经存在
	if _, ok := m.Apis[msgId]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}
	// 添加msgId与api之间的绑定关系
	m.Apis[msgId] = router
	log.Info("Add api msgId = ", msgId)
}

// 以非阻塞方式处理消息
func (m *MsgHandle) DoMsgHandler(request *Request) {
	handler, ok := m.Apis[request.GetMsgID()]
	if !ok {
		log.Error("api msgid = ", request.GetMsgID(), " is not found")
		return
	}

	// handler.PreHandle(request)

	handler.Handle(request)

	// handler.PostHandle(request)
}

// 开启一个Worker goroutine
func (m *MsgHandle) StartOneWorker(workerID int, taskQueue chan *Request) {
	log.Info("WorkerID = ", workerID, " is started.")
	// 当前worker接收 chan Request
	for {
		select {
		case request := <-taskQueue:
			m.DoMsgHandler(request)
		}
	}
}

// 启动worker工作池
func (m *MsgHandle) StartWorkerPool() {
	// 遍例需要启动的worker数量，依次启动
	for i := 0; i < int(m.WorkerPoolSize); i++ {
		// 一个worker启动，给当前worker对应的任务队列开辟空间
		m.TaskQueue[i] = make(chan *Request, globalobj.GlobalObject.MaxWorkerTaskLen)
		go m.StartOneWorker(i, m.TaskQueue[i])
	}
}

// 发送request进入chan Request
func (m *MsgHandle) SendMsgToTaskQueue(request *Request) {

	// 取余分配
	workerID := request.GetConnection().ConnID % m.WorkerPoolSize
	log.Info("add ConnID=", request.GetConnection().GetConnID(), "request msgID=", request.GetMsgID(), " to workerID=", workerID)

	m.TaskQueue[workerID] <- request
}
