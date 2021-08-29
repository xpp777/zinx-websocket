package server

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/xiaomingping/zinx-websocket/configs"
	"github.com/xiaomingping/zinx-websocket/iserverface"
	"github.com/xiaomingping/zinx-websocket/utils"
	"net/http"
)

var (
	Upgrader = websocket.Upgrader{
		// 允许所有CORS跨域请求
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	cfg      configs.WebSocketConfig
	GWServer iserverface.IServer
)

type Server struct {
	// Server的消息管理模块
	MsgHandler iserverface.IMsgHandle
	// 当前Server链接管理器
	ConnMgr iserverface.IConnMgr
	// 当前Server连接创建时的hook函数
	OnConnStart func(conn iserverface.IConnection)
	// 当前Server连接断开时的hook函数
	OnConnStop func(conn iserverface.IConnection)
	packet     iserverface.Packet
}

func NewServer(Cfg configs.WebSocketConfig) iserverface.IServer {
	cfg = Cfg
	s := &Server{
		ConnMgr:    NewConnManager(),
		MsgHandler: NewMsgHandle(),
		packet:     NewPack(),
	}
	GWServer = s
	return s
}

func (s *Server) Start(c *gin.Context) {
	connId := utils.GenUUID()
	wsSocket, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	if s.ConnMgr.Len() >= cfg.MaxConn {
		wsSocket.Close()
		return
	}
	dealConn := NewConnection(s, wsSocket, connId, s.MsgHandler)
	dealConn.Start()
}

// 停止服务
func (s *Server) Stop() {
	s.ConnMgr.ClearConn()
}

// 运行服务
func (s *Server) Serve(c *gin.Context) {
	s.Start(c)
	select {}
}

// 路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *Server) AddRouter(msgId uint32, router iserverface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
}

// 得到链接管理
func (s *Server) GetConnMgr() iserverface.IConnMgr {
	return s.ConnMgr
}

// 设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func(iserverface.IConnection)) {
	s.OnConnStart = hookFunc
}

// 设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func(iserverface.IConnection)) {
	s.OnConnStop = hookFunc
}

// 调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn iserverface.IConnection) {
	if s.OnConnStart != nil {
		logx.Infof("CallOnConnStart.....")
		s.OnConnStart(conn)
	}
}

// 调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn iserverface.IConnection) {
	if s.OnConnStop != nil {
		logx.Infof("CallOnConnStop.....")
		s.OnConnStop(conn)
	}
}

func (s *Server) Packet() iserverface.Packet {
	return s.packet
}
