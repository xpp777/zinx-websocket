package configs

type WebSocketConfig struct {
	Name             string `json:"name,default=websocket"`        // 当前服务名称
	ListenOn         string `json:"listenOn,default=0.0.0.0:8090"` // 当前服务器主机IP和端口号
	MaxConn          int    `json:"maxConn,default=10000"`         // 当前服务器主机允许的最大链接个数
	WorkerPoolSize   uint32 `json:"workerPoolSize,default=10"`     // 业务工作Worker池的数量
	MaxWorkerTaskLen uint32 `json:"maxWorkerTaskLen,default=200"`  // 业务工作Worker对应负责的任务队列最大任务存储数量
	MessageType      int    `json:"messageType,default=1"`         // 消息类型
	HeartBeatTime    int    `json:"heartBeatTime,default=30"`      // 心跳检测时间
}
