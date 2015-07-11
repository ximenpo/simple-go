package simple

import (
	"errors"
	"github.com/garyburd/redigo/redis"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	load_from_none = iota
	load_from_file
	load_from_http
	load_from_redis
)

func getDataSourceType(source string) (loadfrom int) {
	loadfrom = load_from_none
	if match, _ := regexp.MatchString("^http:.*", source); match {
		if match, _ := regexp.MatchString("^http://.+", source); match {
			return load_from_http
		}
	} else if match, _ := regexp.MatchString("^redis:.*", source); match {
		re, _ := regexp.Compile(`^redis://([^/]+)/(.+)$`)
		if re.MatchString(source) {
			if r := re.FindAllStringSubmatch(source, -1); r != nil {
				return load_from_redis
			}
		}
	} else if len(source) != 0 {
		if _, err := os.Stat(source); err == nil || os.IsExist(err) {
			return load_from_file
		}
	}
	return load_from_none
}

func LoadData_FromHttp(url string) (data []byte, err error) {
	r, err := http.Get(url)
	if err != nil {
		return
	}
	defer r.Body.Close()

	return ioutil.ReadAll(r.Body)
}

func LoadData_FromRedis(addr string) (data []byte, err error) {
	var redis_addr, redis_key string
	{
		re, _ := regexp.Compile(`^redis://([^/]+)/(.+)$`)
		if !re.MatchString(addr) {
			return nil, errors.New("invalid redis address")
		}

		if r := re.FindAllStringSubmatch(addr, -1); r == nil {
			return nil, errors.New("malformed redis address")
		} else {
			redis_addr = r[0][1]
			redis_key = r[0][2]
		}
	}

	// parse db/key
	var db int
	var key string
	{
		re, _ := regexp.Compile("^[0-9][0-9]?/(.+)$")
		if re.MatchString(redis_key) {
			r := re.FindAllStringSubmatch(redis_key, -1)
			if r == nil {
				return nil, errors.New("invalid redis source param")
			}
			db, _ = strconv.Atoi(r[0][0])
			key = strings.TrimSpace(r[0][1])
		} else {
			db = 0
			key = strings.TrimSpace(redis_key)
		}
	}

	c, err := redis.Dial("tcp", redis_addr)
	if err != nil {
		return
	}
	defer c.Close()

	// select db
	if _, err = c.Do("SELECT", db); err != nil {
		return
	}

	// get & run lua script
	return redis.Bytes(c.Do("GET", key))
}

func LoadData_FromFile(path string) (data []byte, err error) {
	return ioutil.ReadFile(path)
}

func LoadData(addr string) (data []byte, err error) {
	switch getDataSourceType(addr) {
	case load_from_http:
		return LoadData_FromHttp(addr)
	case load_from_redis:
		return LoadData_FromRedis(addr)
	case load_from_file:
		return LoadData_FromFile(addr)
	}

	return nil, errors.New("invalid source format" + addr)
}
