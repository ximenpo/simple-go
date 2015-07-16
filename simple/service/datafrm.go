package service

import (
	"encoding/binary"
	"errors"
	"github.com/ximenpo/simple-go/simple/connmgr"
	"github.com/ximenpo/simple-go/simple/databuf"
)

//
// 幀定义
//
//  |--body_len(2byte)--|----body----|
//
type DataFrame struct {
	Buf databuf.Buffer
}

func (self *DataFrame) FrameData() []byte {
	return self.Buf.Data()
}

// 读取幀数据
type DataFrameReader struct {
}

func (self *DataFrameReader) ReadFrame(conn *connmgr.Conn) (frame connmgr.Frame, err error) {
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

	frm := &DataFrame{}
	frm.Buf.Assign(body)
	return frm, nil
}

// 写入幀数据
type DataFrameWriter struct {
}

func (self *DataFrameWriter) WriteFrame(conn *connmgr.Conn, frame connmgr.Frame) (err error) {
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
