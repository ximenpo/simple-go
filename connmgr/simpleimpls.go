package connmgr

// 连接配置信息
type ConnConfig struct {
	ReadTimeout  uint // seconds
	WriteTimeout uint // seconds
}

// 读取幀数据
type ConnFrameReader struct {
}

func (self *ConnFrameReader) ReadFrame(conn *Conn, frame *Frame) (err error) {
	// TODO:
	return
}

// 写入幀数据
type ConnFrameWriter struct {
}

func (self *ConnFrameWriter) WriteFrame(conn *Conn, frame *Frame) (err error) {
	// TODO:
	return
}

// 读实现
type ConnReader struct {
}

func (self *ConnReader) ReadLoop(conn *Conn) (err error) {
	// TODO:
	return
}

// 写实现
type ConnWriter struct {
}

func (self *ConnWriter) WriteLoop(conn *Conn) (err error) {
	// TODO:
	return
}

// 写分发
type ConnWriteDispatcher struct {
}

func (self *ConnWriteDispatcher) DispatchLoop() (err error) {
	// TODO:
	return
}
