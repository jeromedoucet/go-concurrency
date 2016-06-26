package client

import (
	"github.com/garyburd/redigo/redis"
	"testing"
)

func Test_reset(t *testing.T) {
	c, _ := redis.Dial("tcp", "192.168.99.100:6379")
	defer c.Close()
	c.Do("DEL", "playerId")
}

