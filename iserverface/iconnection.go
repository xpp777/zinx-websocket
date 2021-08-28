package iserverface

import (
	"github.com/gorilla/websocket"
	"net"
)

type IConnection interface {
	// 启动链接开始工作
	Start()
	// 关闭链接停止工作
	Close()
	// 获取websocket链接
	GetConnection() *websocket.Conn
	// 获取当前连接ID
	GetConnID() string
	// 获取远程客户端地址信息
	RemoteAddr() net.Addr
	// 发送数据
	SendMessage(msgData []byte) error

	IsAlive() bool // 心跳检测
	KeepAlive()    // 更新心跳
}
