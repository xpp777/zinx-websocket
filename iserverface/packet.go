package iserverface

/*
	封包数据和拆包数据
*/
type Packet interface {
	Pack(msg IMessage) ([]byte, error) // 封包方法
	Unpack([]byte) (IMessage, error)   // 拆包方法
}
