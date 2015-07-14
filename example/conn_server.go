package main

import (
	. "github.com/ximenpo/simple-go/simple/connmgr"
	"log"
	"net"
	"time"
)

func main() {

	evt_queue := make(chan *Event, 10)
	var srv = ConnAcceptor{
		Factory: &ConnFactory{
			ReadQueue: evt_queue,
		},
		Handler: &ConnHandler{
			Reader: &ConnReader{
				FrameReader: &SimpleFrameReader{},
			},
			Writer: &ConnWriter{
				FrameWriter: &SimpleFrameWriter{},
			},
		},
	}

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("server started")

	sig_stop := make(chan bool)
	go srv.AcceptLoop(listener, sig_stop)

	var timeout time.Duration = 10
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
						log.Println(evt.Frame.FrameData().Dump())
					}

					if evt.Type == MESSAGE {
						evt.Conn.WriteQueue <- evt
					}
				}
			}
		case <-time.After(time.Second * timeout):
			{
				// no operation, exit
				select {
				case _, ok := <-sig_stop:
					if !ok {
						close(evt_queue)
					}
				default:
					close(sig_stop)
					timeout = 3
				}
			}
		}

		if stopping {
			break
		}
	}

	log.Println("server stopped")
}
