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

func (r *RedisPublisher) Close() {
	r.Conn.Close()
}

func (r *RedisPublisher) Connect(addr string) (err error) {
	r.Conn, err = redis.Dial("tcp", addr)
	return
}

func (r *RedisPublisher) Publish(channel string, value interface{}) (err error) {
	_, err = r.Conn.Do("PUBLISH", channel, value)
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

func (r *RedisSubscriber) Close() {
	r.PubSubConn.Close()
}

func (r *RedisSubscriber) Connect(addr string) (err error) {
	r.PubSubConn.Conn, err = redis.Dial("tcp", addr)
	return
}
