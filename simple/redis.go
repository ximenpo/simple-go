package simple

import (
	"github.com/garyburd/redigo/redis"
)

//
//	redis 发布
//
type RedisPublisher struct {
	redis.Conn
}

func NewRedisPublisher() *RedisPublisher {
	return &RedisPublisher{}
}

func (self *RedisPublisher) Close() {
	self.Conn.Close()
}

func (self *RedisPublisher) Connect(addr string) (err error) {
	self.Conn, err = redis.Dial("tcp", addr)
	return
}

func (self *RedisPublisher) Publish(channel string, value interface{}) (err error) {
	_, err = self.Conn.Do("PUBLISH", channel, value)
	return
}

//
//	redis 订阅
//
type RedisSubscriber struct {
	redis.PubSubConn
}

func NewRedisSubscriber() *RedisSubscriber {
	return &RedisSubscriber{}
}

func (self *RedisSubscriber) Close() {
	self.PubSubConn.Close()
}

func (self *RedisSubscriber) Connect(addr string) (err error) {
	self.PubSubConn.Conn, err = redis.Dial("tcp", addr)
	return
}
