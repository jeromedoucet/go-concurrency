package client_test

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"go-concurrency/drunker/client"
	"go-concurrency/drunker/message"
	"testing"
	"time"
)

type mockRedisC struct {
	t         *testing.T
	countSet  int
	countGet  int
	setVal    interface{}
	orderChan chan *interface{}
}

type mockNsq struct {
	countProduce int
	Val          interface{}
}

func newMock() (mr *mockRedisC, mn *mockNsq) {
	mr = new(mockRedisC)
	mr.orderChan = make(chan *interface{})
	mn = new(mockNsq)
	return
}

func (m *mockRedisC) Set(key string, value interface{}, ttl time.Duration) error {
	m.orderChan <- &value
	m.setVal = value
	m.countSet++
	return nil
}

func (m *mockRedisC) Get(key string) (struct{}, error) {
	m.countGet++
	return *new(struct{}), nil
}

func (m *mockRedisC) Remove(key string) (err error) {
	return
}

func (m *mockNsq) Publish(topic string, body []byte) error {
	m.countProduce++
	m.Val = body
	return nil
}

// test the client with only one producer
func TestSaveOrderInRedis(t *testing.T) {
	mockRedis, mockNsq := newMock()
	c, err := client.StartClient(mockRedis, mockNsq, "myTopic", 1)
	if err != nil {
		t.Errorf("An error has occured during the client starting : %f", err)
	}
	// wait for first order
	<-mockRedis.orderChan
	// and then stop immediatly
	c.StopClient()
	if mockRedis.countSet < 1 && mockNsq.countProduce < 1 {
		t.Errorf("set redis not called")
	} else {
		redisVal := umarshallMess(mockRedis.setVal)
		brokerVal := umarshallMess(mockNsq.Val)
		if redisVal != brokerVal {
			t.Errorf("Value store in Redis %v and value send to broker %v are differents", redisVal, brokerVal)
		}
	}
}

func umarshallMess(data interface{}) message.Order {
	var m message.Order
	b := &bytes.Buffer{}
	binary.Write(b, binary.BigEndian, data)
	json.Unmarshal(b.Bytes(), &m)
	return m
}
