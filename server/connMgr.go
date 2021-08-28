package server

import (
	"errors"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/xiaomingping/zinx-websocket/iserverface"
	"sync"
	"time"
)

/*
	连接管理模块
*/
type ConnManager struct {
	connections map[string]iserverface.IConnection // 管理的连接信息
	connLock    sync.RWMutex                       // 读写连接的读写锁
}

/*
	创建一个链接管理
*/
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[string]iserverface.IConnection),
	}
}

// 添加链接
func (connMgr *ConnManager) Add(conn iserverface.IConnection) {
	// 保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	connMgr.connections[conn.GetConnID()] = conn
	logx.Infof("conn add num = %d", connMgr.Len())
}

// 删除连接
func (connMgr *ConnManager) Remove(conn iserverface.IConnection) {
	// 保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	// 删除连接信息
	delete(connMgr.connections, conn.GetConnID())
	logx.Infof("conn Remove connId = %s  num = %d",conn.GetConnID(),connMgr.Len())
}

// 利用ConnID获取链接
func (connMgr *ConnManager) Get(connID string) (iserverface.IConnection, error) {
	// 保护共享资源Map 加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()
	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

// 获取当前连接
func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

// 清除并停止所有连接
func (connMgr *ConnManager) ClearConn() {
	// 保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	// 停止并删除全部的连接信息
	for connID, conn := range connMgr.connections {
		// 停止
		conn.Close()
		// 删除
		delete(connMgr.connections, connID)
	}
}

func (connMgr *ConnManager) PushAll(msgData []byte) {
	for _, conn := range connMgr.connections {
		conn.SendMessage(msgData)
	}
}

func (connMgr *ConnManager) KeepAlive() {
	var (
		timer *time.Timer
	)
	timer = time.NewTimer(time.Duration(cfg.HeartBeatTime) * time.Second)
	for {
		select {
		case <-timer.C:
			connMgr.keepAlive()
			timer.Reset(time.Duration(cfg.HeartBeatTime) * time.Second)
		}
	}
}

func (connMgr *ConnManager) keepAlive() {
	for _, conn := range connMgr.connections {
		if !conn.IsAlive() {
			conn.Close()
		}
	}
}
