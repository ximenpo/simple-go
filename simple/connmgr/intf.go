package connmgr

import (
	"net"
)

//
//	连接管理模型
//	-	每条连接起两个GO程，一个读，一个写
//	-	连接管理负责Event和Message的汇集和分发
//	-	会有一个或多个写分发GO程，将输出数据包分发到对应连接的写GO程
//
//	事件流程
//	-	连接的读循环将事件抛入连接管理
//	-	处理程序从连接管理读入事件，进行处理，然后再PushEvent或PushMessage
//	-	写分发GO程将需要分发的数据包分发到对应连接的写GO程
//	-	连接写GO程将数据包发送到客户端
//

// Conn工厂
type Factory interface {
	NewConn(conn net.Conn) (ret *Conn, err error) // 生成Conn对象
}

// 读处理循环，可被调用于GO程
type Reader interface {
	ReadLoop(conn *Conn, stop <-chan bool) error // 处理连接读循环
}

// 写处理循环，可被调用于GO程
type Writer interface {
	WriteLoop(conn *Conn, stop <-chan bool) error // 处理连接写循环
}

// 连接处理入口
type Handler interface {
	HandleConn(conn *Conn, stop <-chan bool) error // 处理新连接
}

// 写队列分发循环，可被调用于GO程
type Dispatcher interface {
	DispatchLoop(queue <-chan *Event, stop <-chan bool) error // 分发输出数据包
}

// 接收循环
type Acceptor interface {
	AcceptLoop(listener *net.TCPListener, stop <-chan bool) error // 处理Accept循环
}
