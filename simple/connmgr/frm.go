package connmgr

type Frame interface {
	FrameData() []byte //  data content
}

type FrameReader interface {
	ReadFrame(conn *Conn) (frame Frame, err error) // 读取一帧数据
}

type FrameWriter interface {
	WriteFrame(conn *Conn, frame Frame) error // 读取一帧数据
}
