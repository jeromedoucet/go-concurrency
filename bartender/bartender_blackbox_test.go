package bartender_test

import (
	"bytes"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/stretchr/testify/assert"
	"github.com/vil-coyote-acme/go-concurrency/bartender"
	"github.com/vil-coyote-acme/go-concurrency/commons"
	"net/http"
	"sync"
	"testing"
)

var (
	redisAddr string               = "192.168.99.100:6379" // todo put the host as var ?
	b         *bartender.Bartender = bartender.NewBartender(redisAddr)
)

func Test_bartender_should_handle_order_and_check_validity_into_redis_with_success(t *testing.T) {
	// given
	order := commons.Order{Id: 1, Quantity: 5, Type: commons.Beer, CallBackUrl: "http://some-callback.com"}
	bd, _ := json.Marshal(order)
	// prepare datum into redis
	c, _ := redis.Dial("tcp", redisAddr)
	defer c.Close()
	c.Do("SET", order.Id, string(bd))
	defer c.Do("DEL", order.Id)
	// start bartender
	startBartender(b)
	// when
	resp, err := http.Post("http://127.0.0.1:4343/orders", "application/json", bytes.NewBuffer(bd))
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode)
	val, _ := c.Do("GET", order.Id)
	var upgOrder commons.Order
	unmErr := commons.UnmarshalOrderFromInterface(val, &upgOrder)
	if unmErr == nil {
		assert.Equal(t, order.Id, upgOrder.Id)
		assert.True(t, upgOrder.Valid)
	} else {
		t.Fatal(unmErr.Error())
	}
}

func Test_bartender_should_handle_order_and_check_validity_into_redis_and_return_403_if_already_validated(t *testing.T) {
	// given
	order := commons.Order{Id: 1, Quantity: 5, Type: commons.Beer, CallBackUrl: "http://some-callback.com", Valid: true}
	bd, _ := json.Marshal(order)
	// prepare datum into redis
	c, _ := redis.Dial("tcp", redisAddr)
	defer c.Close()
	c.Do("SET", order.Id, string(bd))
	defer c.Do("DEL", order.Id)
	// start bartender
	startBartender(b)
	// when
	resp, err := http.Post("http://127.0.0.1:4343/orders", "application/json", bytes.NewBuffer(bd))
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 403, resp.StatusCode)
	exist, errExist := c.Do("exists", order.Id)
	assert.Nil(t, errExist)
	assert.Equal(t, int64(1), exist.(int64))
}

func Test_bartender_should_handle_order_and_check_validity_into_redis_and_return_404_if_not_exist(t *testing.T) {
	// given
	order := commons.Order{Id: 1, Quantity: 5, Type: commons.Beer, CallBackUrl: "http://some-callback.com"}
	bd, _ := json.Marshal(order)
	// prepare datum into redis
	c, _ := redis.Dial("tcp", redisAddr)
	defer c.Close()
	// start bartender
	startBartender(b)
	// when
	resp, err := http.Post("http://127.0.0.1:4343/orders", "application/json", bytes.NewBuffer(bd))
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 404, resp.StatusCode)
	exist, errExist := c.Do("exists", order.Id)
	assert.Nil(t, errExist)
	assert.Equal(t, int64(0), exist.(int64))
}

func startBartender(b *bartender.Bartender) {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	go func() {
		wg.Done()
		b.Start()
	}()
	wg.Wait()
}
