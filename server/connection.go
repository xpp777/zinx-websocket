package server

import (
	"errors"
	"github.com/gorilla/websocket"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/xiaomingping/zinx-websocket/iserverface"
	"net"
	"sync"
	"time"
)

type Connection struct {
	Server            iserverface.IServer
	Conn              *websocket.Conn
	connId            string
	outChan           chan []byte
	isClosed          bool
	closeChan         chan byte
	MsgHandle         iserverface.IMsgHandle
	lastHeartBeatTime time.Time
	mutex             sync.Mutex
}

// 初始化链接服务
func NewConnection(server iserverface.IServer, wsSocket *websocket.Conn, connId string, msgHandler iserverface.IMsgHandle) *Connection {
	c := &Connection{
		Server:            server,
		Conn:              wsSocket,
		connId:            connId,
		MsgHandle:         msgHandler,
		outChan:           make(chan []byte, 20),
		closeChan:         make(chan byte),
		lastHeartBeatTime: time.Now(),
	}
	c.Server.GetConnMgr().Add(c)
	return c
}

// 开始
func (c *Connection) Start() {
	go c.readLoop()
	go c.writeLoop()
	c.Server.CallOnConnStart(c)
}

// 关闭连接
func (c *Connection) Close() {
	c.Server.CallOnConnStop(c)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Conn.Close()
	if !c.isClosed {
		c.isClosed = true
		close(c.closeChan)
	}
	c.Server.GetConnMgr().Remove(c)
}

// 获取链接对象
func (c *Connection) GetConnection() *websocket.Conn {
	return c.Conn
}

// 获取链接ID
func (c *Connection) GetConnID() string {
	return c.connId
}

// 获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

// 读websocket
func (c *Connection) readLoop() {
	for {
		messageType, msgData, err := c.Conn.ReadMessage()
		if err != nil {
			if messageType == -1 || messageType == websocket.CloseMessage {
				goto ERR
			}
			logx.Error(err)
			goto ERR
		}
		// 拆包，得到msgID 和 data 放在msg中
		message, err := c.Server.Packet().Unpack(msgData)
		if err != nil {
			logx.Error("unpack error ", err)
			goto ERR
		}
		// 得到当前客户端请求的Request数据
		req := Request{
			conn: c,
			msg:  message,
		}
		if cfg.WorkerPoolSize > 0 {
			// 已经启动工作池机制，将消息交给Worker处理
			c.MsgHandle.SendMsgToTaskQueue(&req)
		} else {
			// 从绑定好的消息和对应的处理方法中执行对应的Handle方法
			go c.MsgHandle.DoMsgHandler(&req)
		}

	}

ERR:
	c.Close()
}

// 写websocket
func (c *Connection) writeLoop() {
	var (
		err error
	)
	for {
		select {
		case message := <-c.outChan:
			c.Conn.SetWriteDeadline(time.Now().Add(time.Second * 10))
			if err = c.Conn.WriteMessage(cfg.MessageType, message); err != nil {
				logx.Error(err)
				goto ERR
			}
		case <-c.closeChan:
			goto CLOSED
		}
	}
ERR:
	c.Close()
CLOSED:
}

// 发送消息
func (c *Connection) SendMessage(msgData []byte) (err error) {
	select {
	case c.outChan <- msgData:
	case <-c.closeChan:
		err = errors.New("ERR_CONNECTION_LOSS")
	default: // 写操作不会阻塞, 因为channel已经预留给websocket一定的缓冲空间
		err = errors.New("ERR_SEND_MESSAGE_FULL")
	}
	return
}

// 检测心跳
func (c *Connection) IsAlive() bool {
	var (
		now = time.Now()
	)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.isClosed || now.Sub(c.lastHeartBeatTime) > time.Duration(cfg.HeartBeatTime)*time.Second {
		return false
	}
	return true

}

// 更新心跳
func (c *Connection) KeepAlive() {
	var (
		now = time.Now()
	)
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.lastHeartBeatTime = now
}
