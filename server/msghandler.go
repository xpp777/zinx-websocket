package server

import (
	"fmt"
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/xiaomingping/zinx-websocket/iserverface"
	"math/rand"
)

type MsgHandle struct {
	Apis           map[uint32]iserverface.IRouter // 存放每个MsgId 所对应的处理方法的map属性
	WorkerPoolSize uint32                         // 业务工作Worker池的数量
	TaskQueue      []chan iserverface.IRequest    // Worker负责取任务的消息队列
}

func NewMsgHandle() *MsgHandle {
	m := &MsgHandle{
		Apis:           make(map[uint32]iserverface.IRouter),
		WorkerPoolSize: cfg.WorkerPoolSize,
		TaskQueue:      make([]chan iserverface.IRequest, cfg.WorkerPoolSize),
	}
	m.StartWorkerPool()
	return m
}

// 将消息交给TaskQueue,由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request iserverface.IRequest) {
	if len(mh.TaskQueue) == 0 {
		return
	}
	i := rand.Intn(len(mh.TaskQueue))
	// 将请求消息发送给任务队列
	mh.TaskQueue[i] <- request
}

// 马上以非阻塞方式处理消息
func (mh *MsgHandle) DoMsgHandler(request iserverface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		logx.Errorf("msg handler msgId = %d is not found!", request.GetMsgID())
		return
	}
	// 执行对应处理方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgId uint32, router iserverface.IRouter) {
	// 1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgId]; ok {
		panic(fmt.Sprintf("repeated api , msgId = %d", msgId))
	}
	// 2 添加msg与api的绑定关系
	mh.Apis[msgId] = router
	logx.Infof("msg handler add msgId = %d", msgId)
}

// 启动一个Worker工作流程
func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan iserverface.IRequest) {
	logx.Infof("msg handler start worker = %d", workerID)
	// 不断的等待队列中的消息
	for {
		select {
		// 有消息则取出队列的Request，并执行绑定的业务方法
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// 启动worker工作池
func (mh *MsgHandle) StartWorkerPool() {
	// 遍历需要启动worker的数量，依此启动
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 一个worker被启动
		// 给当前worker对应的任务队列开辟空间
		mh.TaskQueue[i] = make(chan iserverface.IRequest, cfg.MaxWorkerTaskLen)
		// 启动当前Worker，阻塞的等待对应的任务队列是否有消息传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}
