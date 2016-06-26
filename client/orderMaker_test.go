package client

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"sync"
	"net/http"
	"github.com/vil-coyote-acme/go-concurrency/commons"
	"github.com/garyburd/redigo/redis"
	"strings"
)

func TestOrderMaker(t *testing.T) {
	url := "http://calback:4444"
	wg := new(sync.WaitGroup)
	wg.Add(1)
	mux := http.NewServeMux()
	mux.HandleFunc("/orders", func(rw http.ResponseWriter, rq *http.Request) {
		assert.Equal(t, http.MethodPost, rq.Method)

		var order commons.Order
		commons.UnmarshalOrderFromHttp(rq, &order)

		c, _ := redis.Dial("tcp", "192.168.99.100:6379")
		defer c.Close()
		exist, _ := c.Do("EXISTS", order.Id)
		assert.Equal(t, int64(1), exist.(int64))
		assert.True(t, strings.HasPrefix(order.CallBackUrl, url + "/bill/playerId/"))
		wg.Done()
		rw.WriteHeader(200)
	})
	srv := httptest.NewServer(mux)
	unReg := make(chan commons.Registration)
	reg := commons.Registration{PlayerId:"playerId", Ip:srv.URL}

	startNewOrderMaker(url, "192.168.99.100:6379", reg, unReg)
	wg.Wait()
}
