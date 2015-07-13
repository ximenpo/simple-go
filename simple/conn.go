package simple

//
//	连接管理模型
//	-	每条连接起两个GO程，一个读，一个写
//	-	连接管理负责Event和Message的汇集和分发
//	-	会有一个或多个写分发GO程，将输出数据包分发到对应连接的写GO程
//
//	事件流程
//	-	连接的读循环将事件抛入连接管理
//	-	处理程序从连接管理读入事件，进行处理，然后再PushEvent或PushMessage
//	-	写分发GO程将需要分发的数据包分发到对应连接的写GO程
//	-	连接写GO程将数据包发送到客户端
//

import (
	buf "github.com/ximenpo/simple-go/simple/databuf"
	"net"
)

//
// 单条连接
//
type Conn struct {
	NetConn      net.Conn // 连接
	Tag          uint32   // 标记
	Disconnected bool     // 是否连接状态
}

func (c *Conn) Close() {
	c.Disconnected = true
	c.NetConn.Close()
}

//	连接处理接口
type ConnHandler interface {
	ReadLoop(params *ConnReadLoopParams) error   // 处理连接读循环
	WriteLoop(params *ConnWriteLoopParams) error // 处理连接写循环
}

//	连接管理接口
type ConnMgrHandler interface {
	PushEvent(evt *ConnEvent)     // 处理 ConnectionEvent
}

//	写分发接口
type WriteDispatcher interface {
	WriteDispatchLoop(queue <-chan *ConnEvent) // 分发输出数据包
}

//////////////////////////////////////////////
//
//	相关类型
//
//////////////////////////////////////////////

//
// 连接事件类型
//
const (
	CONN_EVENT_NONE   = iota // 无消息
	CONN_CONNECTED           // 已连接
	CONN_MESSAGE             // 接收到数据包
	CONN_DISCONNECTED        // 连接断开
	CONN_EVENT_CUSTOM = 0x10 // 自定义事件，必须>=EVENT_CUSTOM
)

//
// 连接事件
//
type ConnEvent struct {
	Type int        //  msg type
	Conn Conn       //	src/target conn
	Data buf.Buffer //  msg data
}

//
// 连接配置信息
//
type ConnConfig struct {
	ReadTimeout  uint // seconds
	WriteTimeout uint // seconds
}

//
// 连接读循环参数
//
type ConnReadLoopParams struct {
	Queue   chan<- *ConnEvent // 消息队列
	SigStop <-chan bool       // 停止消息
	Cfg     *ConnConfig       // 配置信息
}

//
// 连接写循环参数
//
type ConnWriteLoopParams struct {
	Queue   <-chan *ConnEvent // 消息队列
	SigStop <-chan bool       // 停止消息
	Cfg     *ConnConfig       // 配置信息
}
