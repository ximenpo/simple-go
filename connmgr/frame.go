package connmgr

import (
	buf "github.com/ximenpo/simple-go/databuf"
)

type Frame struct {
	Type int        // data type
	Data buf.Buffer //  data content
}

type FrameReader interface {
	ReadFrame(conn *Conn, frame *Frame) error // 读取一帧数据
}

type FrameWriter interface {
	WriteFrame(conn *Conn, frame *Frame) error // 读取一帧数据
}
