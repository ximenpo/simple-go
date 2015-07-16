package main

import (
	. "github.com/ximenpo/simple-go/simple/connmgr"
	. "github.com/ximenpo/simple-go/simple/databuf"
	"log"
	"net"
	"time"
)

func main() {

	evt_queue := make(chan *Event, 10)
	var client = &ConnHandler{
		Reader: &ConnReader{
			FrameReader: &SimpleFrameReader{},
		},
		Writer: &ConnWriter{
			FrameWriter: &SimpleFrameWriter{},
		},
	}
	var factory = &ConnFactory{
		ReadQueue: evt_queue,
	}

	net_conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("client started")

	sig_stop := make(chan bool)

	conn, err := factory.NewConn(net_conn)
	if err != nil {
		log.Fatal(err)
	}
	go client.HandleConn(conn, sig_stop)

	var stopping bool
	for {
		select {
		case evt, ok := <-evt_queue:
			{
				if !ok {
					stopping = true
				} else {
					log.Println(evt.Conn, evt.Type)
					if evt.Frame != nil {
						log.Println(evt.Frame.FrameData())
					}

					if evt.Type == DISCONNECTED {
						close(sig_stop)
						close(evt_queue)
					}
				}
			}
		case <-time.After(time.Second * 5):
			{
				var buf Buffer
				NewDataWriter(&buf).Write("0123456789")

				conn.WriteQueue <- &Event{
					MESSAGE,
					conn,
					&SimpleFrame{buf},
				}
			}
		}

		if stopping {
			break
		}
	}

	log.Println("client stopped")
}
