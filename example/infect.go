package main

import (
    . "github.com/ximenpo/simple-go/simple/connmgr"
    "github.com/ximenpo/simple-go/simple"
    "log"
)

func main() {
	evt_queue := make(chan *Event, 10)
    var client = ConnClient{
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

    log.Println(client)
    log.Println(client.Factory)
    log.Println(client.Handler)
    log.Println("------------------------")

    var   reader Reader = &ConnReader{
            FrameReader: &SimpleFrameReader{},
    }
    simple.Infect_InterfaceField(&client, reader)

    log.Println(client)
    log.Println(client.Factory)
    log.Println(client.Handler)
    log.Println("------------------------")

    var config = ConnConfig{10, 20, 30}
    simple.Infect_InterfaceField(&client, &config)

    log.Println(client)
    log.Println(client.Factory)
    log.Println(client.Handler)
}
