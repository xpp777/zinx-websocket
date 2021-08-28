package iserverface

/*
	将请求的一个消息封装到message中，定义抽象层接口
*/
type IMessage interface {
	GetMsgID() uint32     // 获取消息ID
	GetData() interface{} // 获取消息内容
}
