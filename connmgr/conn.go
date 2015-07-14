package connmgr

import (
	"net"
)

// 单条连接
type Conn struct {
	NetConn      net.Conn // 连接
	Tag          uint32   // 标记
	Disconnected bool     // 是否连接状态
}

func (self *Conn) Close() {
	self.Disconnected = true
	self.NetConn.Close()
}
