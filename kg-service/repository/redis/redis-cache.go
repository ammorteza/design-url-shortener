package redis

import (
	"errors"
	"github.com/ammorteza/clean_architecture/repository"
	"github.com/gomodule/redigo/redis"
	"log"
)

type Cache struct {
	connPool 		*redis.Pool
}

func NewRedisPool() repository.Cache{
	return &Cache{
		connPool: &redis.Pool{
			MaxIdle: 10,
			MaxActive: 12000,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", "172.28.1.112:6379")
				if err != nil{
					panic(err.Error())
				}
				return c, err
			},
		},
	}
}

func (c *Cache)getConn() redis.Conn{
	return c.connPool.Get()
}

func (c *Cache)Ping() (err error){
	conn := c.getConn()
	defer func() {
		err = conn.Close()
	}()

	s, err := redis.String(conn.(redis.Conn).Do("PING"))
	if err != nil{
		return err
	}
	log.Printf("redis PING response = %s\n", s)
	return
}

func (c *Cache)Set(key, value string) (err error){
	conn := c.getConn()
	defer func() {
		err = conn.Close()
	}()
	_, err = conn.(redis.Conn).Do("SET", key, value)
	if err != nil{
		return err
	}
	return
}

func (c *Cache)Get(key string) (val string, err error){
	conn := c.getConn()
	defer func() {
		err = conn.Close()
	}()

	val, err = redis.String(conn.(redis.Conn).Do("GET", key))
	return
}

func (c *Cache)Remove(key string) (err error){
	conn := c.getConn()
	defer func() {
		err = conn.Close()
	}()

	_, err = conn.Do("DEL", key)
	return
}

func (c *Cache)Push(listName, val string) (err error){
	conn := c.getConn()
	defer func() {
		err = conn.Close()
	}()

	_, err = conn.Do("RPUSH" , listName, val)
	return
}

func (c *Cache)Pop(listName string) (string, error){
	var err error
	var val string
	var count int64

	conn := c.getConn()
	defer func() {
		err = conn.Close()
	}()

	count, err = c.ListLen(listName)
	if err != nil || count == 0{
		err = errors.New("cannot find any unique key")
	}else{
		var temp interface{}
		temp, err = conn.Do("RPOP", listName)
		val = string(temp.([]uint8))
	}

	return val, err
}

func (c *Cache)ListLen(listName string) (count int64, err error){
	conn := c.getConn()
	defer func() {
		err = conn.Close()
	}()

	temp, err := conn.Do("LLEN", listName)
	if err != nil{
		return 0, err
	}
	count = temp.(int64)
	return
}