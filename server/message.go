package server

type Message struct {
	MsgID uint32      `json:"msgId"` // 业务消息ID
	Data  interface{} `json:"data"`  // 消息的内容
}

// 创建一个Message消息包
func NewMsg(msgID uint32, data interface{}) *Message {
	return &Message{
		MsgID: msgID,
		Data:  data,
	}
}

// 获取消息类型
func (msg *Message) GetMsgID() uint32 {
	return msg.MsgID
}

// 获取消息内容
func (msg *Message) GetData() interface{} {
	return msg.Data
}
