package simple

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

func (self *Conn) Close() {
	self.Disconnected = true
	self.NetConn.Close()
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
