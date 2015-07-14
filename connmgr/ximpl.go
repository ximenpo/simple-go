package connmgr

import (
	"errors"
	"fmt"
	"time"
)

// 连接配置信息
type ConnConfig struct {
	ReadTimeout  uint // seconds
	WriteTimeout uint // seconds
}

// 读实现
type ConnReader struct {
	Config      *ConnConfig // optional
	FrameReader FrameReader // must be set
}

func (self *ConnReader) ReadLoop(conn *Conn, stop <-chan bool) (err error) {
	if conn == nil {
		return errors.New("conn must not be nil")
	}
	if conn.ReadQueue == nil {
		return errors.New("conn read queue must not be nil")
	}
	if self.FrameReader == nil {
		return errors.New("frame reader must not be nil")
	}

	for {
		var evt = new(Event)

		cfg := self.Config
		if cfg != nil && cfg.ReadTimeout >= 1 {
			d := time.Second * time.Duration(cfg.ReadTimeout)
			conn.NetConn.SetReadDeadline(time.Now().Add(d))
		}

		if evt.Frame, err = self.FrameReader.ReadFrame(conn); err != nil {
			return
		}

		if cfg != nil {
			conn.NetConn.SetReadDeadline(time.Time{})
		}

		if evt.Frame == nil {
			return errors.New("ReadFrame returns nil Frame")
		}

		if data := evt.Frame.FrameData(); data != nil {
			data.Rewind()
		}

		conn.ReadQueue <- evt
	}
	return
}

// 写实现
type ConnWriter struct {
	Config      *ConnConfig // optional
	FrameWriter FrameWriter // must be set
}

func (self *ConnWriter) WriteLoop(conn *Conn, stop <-chan bool) (err error) {
	if conn == nil {
		return errors.New("conn must not be nil")
	}
	if conn.WriteQueue == nil {
		return errors.New("conn write queue must not be nil")
	}
	if self.FrameWriter == nil {
		return errors.New("frame writer must not be nil")
	}

	for {
		select {
		case evt, ok := <-conn.WriteQueue:
			if !ok {
				return errors.New("queue closed")
			}

			cfg := self.Config
			if cfg != nil && cfg.WriteTimeout >= 1 {
				d := time.Second * time.Duration(cfg.WriteTimeout)
				conn.NetConn.SetWriteDeadline(time.Now().Add(d))
			}

			if evt.MsgType == MESSAGE {
				if err = self.FrameWriter.WriteFrame(conn, evt.Frame); err != nil {
					return
				}
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

// 连接处理入口
type ConnHandler struct {
	Reader Reader
	Writer Writer
	Stop   chan bool
}

func (self *ConnHandler) HandleConn(conn *Conn) (err error) {
	if conn == nil {
		return errors.New("conn must not be nil")
	}
	if conn.ReadQueue == nil {
		return errors.New("conn read queue must not be nil")
	}
	if self.Reader == nil {
		return errors.New("Reader must not be nil")
	}
	if self.Writer == nil {
		return errors.New("Writer must not be nil")
	}

	var sigstop = self.Stop
	if sigstop == nil {
		sigstop = make(chan bool, 0)
	}

	// connected event
	conn.ReadQueue <- &Event{CONNECTED, conn, nil}
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
		if sigstop != nil {
			close(sigstop)
		}

		// disconnect event
		conn.ReadQueue <- &Event{CONNECTED, conn, nil}
		conn.Close()
	}()

	// writer loop
	go self.Writer.WriteLoop(conn, sigstop)

	// reader loop
	if err = self.Reader.ReadLoop(conn, sigstop); err != nil {
		return
	}

	return
}

// 写分发
type ConnWriteDispatcher struct {
}

func (self *ConnWriteDispatcher) DispatchLoop(queue <-chan *Event) (err error) {
	for {
		select {
		case evt, ok := <-queue:
			if !ok {
				return errors.New("queue closed")
			}

			if !evt.Conn.Disconnected {
				evt.Conn.WriteQueue <- evt
			} else {
				// drop it
			}
		}
	}
	return
}
