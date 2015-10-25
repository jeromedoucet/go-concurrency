package database
import (
	"time"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"log"
)

type Redis struct {
	con redis.Conn
}

func (m * Redis) Set(key string, value interface{}, ttl time.Duration) (err error)  {
	_, e := m.con.Do("SETEX", key, strconv.Itoa(int(ttl)), value)
	if e != nil {
		err = e
	}
	return
}

func (m * Redis) Get(key string) (val struct{}, err error) {
	v, e := m.con.Do("GET", key)
	if e != nil {
		err = e
	} else {
		val = v.(struct {})
	}
	return
}

func NewRedis(addr string) (r *Redis, err error) {
	log.Println("try to connect to redis on " + addr)
	c, e := redis.Dial("tcp", addr)
	if e!= nil {
		err = e
		return
	}
	r = new(Redis)
	r.con = c
	log.Println("Connection to redis successFull")
	return
}
