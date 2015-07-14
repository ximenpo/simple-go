package connmgr

// 事件类型
const (
	EVENT_NONE   = iota // 无消息
	CONNECTED           // 已连接
	MESSAGE             // 接收到数据包
	DISCONNECTED        // 连接断开
	EVENT_CUSTOM = 0x10 // 自定义事件，必须>=EVENT_CUSTOM
)

// 事件
type Event struct {
	Type  int   //  msg type
	Conn  *Conn //  src/target conn
	Frame Frame //  msg data
}
