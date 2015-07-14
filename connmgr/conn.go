package connmgr

import (
	"net"
)

// 单条连接
type Conn struct {
	NetConn      net.Conn // 连接
	Tag          uint32   // 标记
	Disconnected bool     // 是否连接状态

	ReadQueue  chan *Event // 读入队列，待处理，由外部设置
	WriteQueue chan *Event // 写出队列，待发送，由外部设置
}

func (self *Conn) Close() {
	self.Disconnected = true
	self.NetConn.Close()
}
