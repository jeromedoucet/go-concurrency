package database

import (
	"github.com/garyburd/redigo/redis"
	"log"
	"strconv"
	"time"
)

type Redis struct {
	con redis.Conn
}

func (m *Redis) Set(key string, value interface{}, ttl time.Duration) (err error) {
	_, e := m.con.Do("SETEX", key, strconv.Itoa(int(ttl)), value)
	if e != nil {
		err = e
	}
	return
}

func (m *Redis) Get(key string) (val interface{}, err error) {
	v, e := m.con.Do("GET", key)
	if e != nil {
		err = e
	} else {
		val = v
	}
	return
}

func (m *Redis) Remove(key string) (err error) {
	_, e := m.con.Do("DEL", key)
	if e != nil {
		err = e
	}
	return
}

func (m *Redis) Close() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("error : %s", r)
		}
	}()
	m.con.Close()
}

func NewRedis(addr string) (r *Redis, err error) {
	c, e := redis.Dial("tcp", addr)
	if e != nil {
		err = e
		return
	}
	r = new(Redis)
	r.con = c
	return
}
