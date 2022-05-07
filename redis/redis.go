package redis

import (
	"encoding/json"
	"fmt"
	redigo "github.com/gomodule/redigo/redis"
	"sync"
	"time"
)

var pool *redigo.Pool
var once sync.Once

//TODO 所有操作实链式

type RedisConnectOption struct {
	Uri  string
	Auth string
	Db   int
}
type RedisRsult struct {
	Reply interface{}
	Err   error
}

func initRedis(args ...string) *redigo.Pool {
	uri := "127.0.0.1:6379"
	auth := ""
	dbNum := 0
	timeout :=  time.Duration( 10) * time.Second
	maxIdle := 5
    maxActive := 10
	pool := &redigo.Pool{
		Dial: func() (conn redigo.Conn, e error) {
			conn, err := redigo.Dial(
				"tcp",
				uri,
				redigo.DialPassword(auth),
				redigo.DialDatabase(dbNum),
				redigo.DialConnectTimeout(timeout),
				redigo.DialReadTimeout(timeout),
				redigo.DialWriteTimeout(timeout),
			)
			if err != nil {
				return nil, fmt.Errorf("redis connection error: %s", err)
			}
			return
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: timeout,
		Wait:        true,
	}
	return pool
}


func Pool(args ...string) *redigo.Pool {
	once.Do(func() {
		pool = initRedis(args...)
	})
	return pool
}


// GET
func Get(key string) (reply interface{}, err error) {
	con := Pool().Get()
	if err := con.Err(); err != nil {
		return nil, err
	}
	defer  con.Close()
	return con.Do("GET", key)
}

// 	SET
func Set(key string, val interface{}, expire int64) error {
	c := Pool().Get()
	defer c.Close()
	value, err := encode(val)

	if err != nil {
		return err
	}
	if expire > 0 {
		_, err = c.Do("SETEX", key, expire, value)

		return err
	}
	_, err = c.Do("SET", key, value)

	return err
}

func encode(val interface{}) (interface{}, error) {
	var value interface{}
	switch v := val.(type) {
	case string, int, uint, int8, int16, int32, int64, float32, float64, bool:
		value = v
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		value = string(b)
	}
	return value, nil
}

// 反序列化保存的值
func decode(reply interface{}, err error, val interface{}) error {
	str, err := String(reply, err)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(str), val)
}

func String(reply interface{}, err error) (string, error) {
	return redigo.String(reply, err)
}
