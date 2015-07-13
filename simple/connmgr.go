package simple

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

//
// 连接读循环参数
//
type ConnReaderParams struct {
	Queue chan<- *ConnEvent // 消息队列
	Stop  <-chan bool       // 停止消息
	Cfg   *ConnConfig       // 配置信息
}

//
// 连接写循环参数
//
type ConnWriterParams struct {
	Queue <-chan *ConnEvent // 消息队列
	Stop  <-chan bool       // 停止消息
	Cfg   *ConnConfig       // 配置信息
}

// 读循环
type ConnReader interface {
	ReaderLoop(conn *Conn, params ConnReaderParams) error // 处理连接读循环
}

// 写循环
type ConnWriter interface {
	WriterLoop(conn *Conn, params ConnWriterParams) error // 处理连接写循环
}

// 写队列分发循环
type ConnWriteDispatcher interface {
	WriteDispatcherLoop(queue <-chan *ConnEvent) // 分发输出数据包
}

//	连接处理接口
type ConnMgrHandler interface {
	//PushEvent(evt *ConnEvent) // 处理 ConnectionEvent
}

type ConnHandler struct {
	Reader          ConnReader
	Writer          ConnWriter
	WriterDispacher ConnWriteDispatcher
}

type DefaultConnHandler struct {
	ConnHandler
}

//func (self *DefaultConnHandler)

func HandleConn(handler *ConnHandler) {

}
