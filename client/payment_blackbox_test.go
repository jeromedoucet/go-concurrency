package client

import (
	"testing"
	"net/http"
	"github.com/vil-coyote-acme/go-concurrency/commons"
	"net/http/httptest"
	"github.com/garyburd/redigo/redis"
	"encoding/json"
	"github.com/stretchr/testify/assert"
)

var (
	rAddr string = "192.168.99.100:6379"
)

func Test_must_perform_a_payment_and_build_notification(t *testing.T) {
	// given
	mux := http.NewServeMux()
	s := httptest.NewUnstartedServer(mux)
	playerId := "playerId"
	order := commons.Order{Id: 1, Quantity: 5, Type: commons.Beer, CallBackUrl: s.URL + "/bill/playerId/1", Valid:true}
	bdOrder, _ := json.Marshal(order)
	credit := commons.Credit{PlayerId:playerId, Score:10}
	bdCredit, _ := json.Marshal(credit)
	c, _ := redis.Dial("tcp", rAddr)
	defer c.Close()
	c.Do("SET", order.Id, string(bdOrder))
	c.Do("SET", credit.PlayerId, string(bdCredit))
	notifChan := make(chan commons.Notification, 100)
	initPaymentHandling(mux, rAddr, notifChan)
	s.Start()
	// when
	resp, err := http.Get(s.URL + "/bill/playerId/1")
	// then
	notif := <-notifChan
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, playerId, notif.PlayerId)
	assert.Equal(t, 15,  notif.Score)
}
