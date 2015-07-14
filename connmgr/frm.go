package connmgr

import (
	buf "github.com/ximenpo/simple-go/databuf"
)

type Frame interface {
	FrameData() *buf.Buffer //  data content
}

type FrameReader interface {
	ReadFrame(conn *Conn) (frame Frame, err error) // 读取一帧数据
}

type FrameWriter interface {
	WriteFrame(conn *Conn, frame Frame) error // 读取一帧数据
}
