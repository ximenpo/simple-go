package connmgr

import (
	"errors"
	"fmt"
	"net"
	"time"
)

// 连接配置信息
type ConnConfig struct {
	ReadTimeout  uint // seconds
	WriteTimeout uint // seconds

	WriteQueueSize uint // output queue size
}

// Conn工厂
type ConnFactory struct {
	Config *ConnConfig // optional

	ReadQueue chan *Event // 读入队列，待处理，由外部设置
}

func (self *ConnFactory) NewConn(net_conn net.Conn) (ret *Conn, err error) {
	if self.ReadQueue == nil {
		return nil, errors.New("read queue must not be nil")
	}

	cfg := self.Config
	queue_size := uint(16)
	if cfg != nil {
		queue_size = cfg.WriteQueueSize
	}

	return &Conn{
		NetConn:      net_conn,
		Disconnected: false,
		ReadQueue:    self.ReadQueue,
		WriteQueue:   make(chan *Event, queue_size),
		Tag:          0,
	}, nil
}

// 读实现
type ConnReader struct {
	Config *ConnConfig // optional

	FrameReader FrameReader // must be set
}

func (self *ConnReader) ReadLoop(conn *Conn, stop <-chan bool) (err error) {
	if conn == nil {
		return errors.New("conn must not be nil")
	}
	if stop == nil {
		return errors.New("stop chan must not be nil")
	}
	if conn.ReadQueue == nil {
		return errors.New("conn read queue must not be nil")
	}
	if self.FrameReader == nil {
		return errors.New("frame reader must not be nil")
	}

	for {
		var evt = &Event{
			Type: MESSAGE,
			Conn: conn,
		}

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

		conn.ReadQueue <- evt
	}
	return
}

// 写实现
type ConnWriter struct {
	Config *ConnConfig // optional

	FrameWriter FrameWriter // must be set
}

func (self *ConnWriter) WriteLoop(conn *Conn, stop <-chan bool) (err error) {
	if conn == nil {
		return errors.New("conn must not be nil")
	}
	if stop == nil {
		return errors.New("stop chan must not be nil")
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

			if evt.Type == MESSAGE {
				if err = self.FrameWriter.WriteFrame(conn, evt.Frame); err != nil {
					return
				}
			}

			if cfg != nil {
				conn.NetConn.SetWriteDeadline(time.Time{})
			}

		case <-stop:
			return
		}
	}

	return
}

// 连接处理入口
type ConnHandler struct {
	Config *ConnConfig // optional

	Reader Reader
	Writer Writer
}

func (self *ConnHandler) HandleConn(conn *Conn, stop <-chan bool) (err error) {
	if conn == nil {
		return errors.New("conn must not be nil")
	}
	if conn.ReadQueue == nil {
		return errors.New("conn read queue must not be nil")
	}
	if stop == nil {
		return errors.New("stop chan must not be nil")
	}
	if self.Reader == nil {
		return errors.New("Reader must not be nil")
	}
	if self.Writer == nil {
		return errors.New("Writer must not be nil")
	}

	// connected event
	conn.ReadQueue <- &Event{CONNECTED, conn, nil}
	defer func() {
		if e := recover(); e != nil {
			err = errors.New(fmt.Sprint(e))
		}
		// disconnect event
		conn.ReadQueue <- &Event{DISCONNECTED, conn, nil}
		conn.Close()
	}()

	// writer loop
	go self.Writer.WriteLoop(conn, stop)

	// reader loop
	if err = self.Reader.ReadLoop(conn, stop); err != nil {
		return
	}

	return
}

// 接收循环
type ConnAcceptor struct {
	Config *ConnConfig // optional

	Factory Factory
	Handler Handler

	SyncHandler bool // false: accept & async process
}

func (self *ConnAcceptor) AcceptLoop(listener net.Listener, stop <-chan bool) (err error) {
	if listener == nil {
		return errors.New("listener must not be nil")
	}
	if stop == nil {
		return errors.New("stop chan must not be nil")
	}
	if self.Handler == nil {
		return errors.New("Handler must not be nil")
	}
	if self.Factory == nil {
		return errors.New("Factory must not be nil")
	}

	var stopping = false
	var net_conn net.Conn
	for {
		net_conn, err = listener.Accept()
		if err != nil {
			return
		}

		if conn, _ := self.Factory.NewConn(net_conn); conn != nil {
			if self.SyncHandler {
				self.Handler.HandleConn(conn, stop)
			} else {
				go self.Handler.HandleConn(conn, stop)
			}
		}

		// check stop
		select {
		case _, ok := <-stop:
			if !ok {
				stopping = true
			}
		default:
			stopping = false
		}

		if stopping {
			break
		}
	}

	return
}

// 写分发
type ConnWriteDispatcher struct {
	Config *ConnConfig // optional
}

func (self *ConnWriteDispatcher) DispatchLoop(queue <-chan *Event, stop <-chan bool) (err error) {
	if queue == nil {
		return errors.New("write dispatch queue must not be nil")
	}
	if stop == nil {
		return errors.New("stop chan must not be nil")
	}

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
		case <-stop:
			return
		}
	}
	return
}
