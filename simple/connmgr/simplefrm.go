package connmgr

import (
	"bytes"
	"encoding/binary"
	"errors"
)

//
// 幀定义
//
//  |--body_len(2byte)--|----body----|
//
type SimpleFrame struct {
	Buf *bytes.Buffer
}

func (self *SimpleFrame) FrameData() []byte {
	if self.Buf == nil {
		return nil
	}
	return self.Buf.Bytes()
}

// 读取幀数据
type SimpleFrameReader struct {
}

func (self *SimpleFrameReader) ReadFrame(conn *Conn) (frame Frame, err error) {
	if conn.Disconnected {
		return nil, errors.New("conn was closed")
	}

	var body_len uint16
	if err = binary.Read(conn.NetConn, binary.BigEndian, &body_len); err != nil {
		return
	}

	var body []byte
	if body_len > 0 {
		body = make([]byte, body_len)
		if err = binary.Read(conn.NetConn, binary.BigEndian, body); err != nil {
			return
		}
	}

	frm := &SimpleFrame{bytes.NewBuffer(body)}
	return frm, nil
}

// 写入幀数据
type SimpleFrameWriter struct {
}

func (self *SimpleFrameWriter) WriteFrame(conn *Conn, frame Frame) (err error) {
	if conn.Disconnected {
		return errors.New("conn was closed")
	}

	var body_len uint16

	var body = frame.FrameData()
	if body != nil {
		body_len = uint16(len(body))
	}

	if err = binary.Write(conn.NetConn, binary.BigEndian, body_len); err != nil {
		return
	}

	if body_len > 0 {
		if err = binary.Write(conn.NetConn, binary.BigEndian, body); err != nil {
			return
		}
	}

	return
}
