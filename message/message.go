package message

type Message struct {
	Id      uint32 // 消息Id
	DataLen uint32 // 消息长度
	Data    []byte // 消息的内容
}

// 创建一个消息
func NewMessage(id uint32, data []byte) *Message {
	return &Message{
		Id:      id,
		DataLen: uint32(len(data)),
		Data:    data,
	}
}

// 获取消息数据长度
func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

// 获取消息ID
func (m *Message) GetMsgId() uint32 {
	return m.Id
}

// 获取消息内容
func (m *Message) GetData() []byte {
	return m.Data
}

// 设置消息数据长度
func (m *Message) SetDataLen(len uint32) {
	m.DataLen = len
}

// 设置消息ID
func (m *Message) SetMsgId(msgId uint32) {
	m.Id = msgId
}

// 设置消息内容
func (m *Message) SetData(data []byte) {
	m.Data = data
}
