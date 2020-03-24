package message

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/hilonfot/network/server/globalobj"
)

type DataPacker interface {
	GetHeadLen() uint32                         // 获取包头长度方法
	Pack(msg *Message) ([]byte, error)          // 封包
	UpPack(binaryData []byte) (*Message, error) // 拆包
}

type DataPack struct{}

func NewDataPack() *DataPack {
	return &DataPack{}
}

func (d *DataPack) GetHeadLen() uint32 {
	// Id uint32(4 byte) + DataLen uint32(4字节)

	return uint32(8)
}

func (d *DataPack) Pack(msg *Message) ([]byte, error) {
	// 创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 写dataLen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	// 写入msgID
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	// 写入data数据
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}
func (d *DataPack) UpPack(binaryData []byte) (*Message, error) {
	// 创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	// 只解压head的信息，得到dataLen 和 msgID
	msg := &Message{}

	// 读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// 读msgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// 判断data的长度是否超出我们允许的最大长度
	if globalobj.GlobalObject.MaxPacketSize > 0 && msg.DataLen > globalobj.GlobalObject.MaxPacketSize {
		return nil, errors.New("Too large msg data received")
	}
	return msg, nil
}

var _ DataPacker = (*DataPack)(nil)
