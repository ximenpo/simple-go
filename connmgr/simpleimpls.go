package connmgr

import (
	"errors"
	"time"
)

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
	Config      *ConnConfig
	Queue       <-chan Event
	FrameReader FrameReader
}

func (self *ConnReader) ReadLoop(conn *Conn, stop <-chan bool) (err error) {
	// TODO:
	return
}

// 写实现
type ConnWriter struct {
	Config      *ConnConfig  // optional
	Queue       <-chan Event // must be set
	FrameWriter FrameWriter  // must be set
}

func (self *ConnWriter) WriteLoop(conn *Conn, stop <-chan bool) (err error) {
	if conn == nil {
		return errors.New("conn must not be nil")
	}
	if self.Queue == nil {
		return errors.New("event queue must not be nil")
	}
	if self.FrameWriter == nil {
		return errors.New("frame writer must not be nil")
	}

	for {
		select {
		case evt, ok := <-self.Queue:
			if !ok {
				return errors.New("queue closed")
			}

			cfg := self.Config
			if cfg != nil && cfg.WriteTimeout >= 1 {
				d := time.Second * time.Duration(cfg.WriteTimeout)
				conn.NetConn.SetWriteDeadline(time.Now().Add(d))
			}

			if err = self.FrameWriter.WriteFrame(conn, &evt.Data); err != nil {
				return
			}

			if cfg != nil {
				conn.NetConn.SetWriteDeadline(time.Time{})
			}

		case _, _ = <-stop:
			return
		}
	}

	return
}

// 写分发
type ConnWriteDispatcher struct {
}

func (self *ConnWriteDispatcher) DispatchLoop() (err error) {
	// TODO:
	return
}
