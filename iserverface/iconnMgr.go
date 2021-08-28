package iserverface

type IConnMgr interface {
	Add(conn IConnection)                   // 添加链接
	Remove(conn IConnection)                // 删除连接
	Get(connID string) (IConnection, error) // 利用ConnID获取链接
	Len() int                               // 获取当前连接
	ClearConn()                             // 删除并停止所有链接
	PushAll(msgData []byte)                 // 广播
	KeepAlive()                             // 心跳维护
}
