package server

import (
	"encoding/json"
	"github.com/xiaomingping/zinx-websocket/iserverface"
)

// Pack 封包拆包类实例，暂时不需要成员
type Pack struct{}

// NewDataPack 封包拆包实例初始化方法
func NewPack() iserverface.Packet {
	return &Pack{}
}

// Pack 封包方法(压缩数据)
func (dp *Pack) Pack(msg iserverface.IMessage) ([]byte, error) {
	return json.Marshal(msg)
}

// Unpack 拆包方法(解压数据)
func (dp *Pack) Unpack(binaryData []byte) (iserverface.IMessage, error) {
	msg := &Message{}
	err := json.Unmarshal(binaryData, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
