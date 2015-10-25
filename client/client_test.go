package client_test
import (
	"testing"
	"gopkg.in/redis.v3"
	"time"
	"go-concurrency/client"
)

type mockRedisC struct {
	t *testing.T
	countSet int
	countGet int
	setVal interface{}
	orderChan chan *interface{}
}

func newMock()(m * mockRedisC){
	m = new(mockRedisC)
	m.orderChan = make(chan *interface{})
	return
}


func (m * mockRedisC) Set(key string, value interface{}, ttl time.Duration) (s *redis.StatusCmd) {
	s = redis.NewStatusCmd()
	m.orderChan <- &value
	m.setVal = value
	m.countSet ++
	return
}

func (m * mockRedisC) Get(key string) (s *redis.StringCmd) {
	s = redis.NewStringCmd()
	m.countGet ++
	return
}

// test the client with only one producer
func TestSaveOrderInRedis(t *testing.T) {
	mock := newMock()
	c, err := client.StartClient(mock, 1)
	if err != nil {
		t.Errorf("An error has occured during the client starting : %f", err)
	}
	<-mock.orderChan
	if mock.countSet < 1 {
		t.Errorf("set redis not called")
	} else {
		t.Log("the last setted value is : ", mock.setVal)
	}
	c.StopClient()
}


