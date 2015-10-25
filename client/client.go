package client
import (
	"time"
	"gopkg.in/redis.v3"
	"go-concurrency/message"
	"go-concurrency/producer"
	"encoding/json"
)

type Client struct {
	mChan chan *message.Order
	stopChan chan struct {}
	redisCl DbClient
	brokerProducer BrokerProducer
	topic string
}

type DbClient interface{
	Set(string, interface{}, time.Duration) *redis.StatusCmd
	Get(string) *redis.StringCmd
}

type BrokerProducer interface {
	Publish(topic string, body []byte) error
}

func StartClient(dbClient DbClient, brokerProducer BrokerProducer, topic string, countP int) (c *Client, err error) {
	c = new(Client)
	// todo mettre le cas d'erreur
	c.redisCl = dbClient
	c.brokerProducer = brokerProducer
	c.mChan = make(chan *message.Order)
	c.stopChan = make(chan struct {}, 1)
	c.topic = topic
	go c.listen()
	for i := 0; i<countP; i++ {
		p := producer.NewProducer(c.stopChan)
		p.Start(c.mChan)
	}
	return
}

func (c * Client) listen() {
	for {
		o := <-c.mChan
		if o != nil {
			json,_ := json.Marshal(o)
			c.redisCl.Set(string(o.Id), json, 10000000)
			c.brokerProducer.Publish(c.topic, json) // TODO do something with err !
		}else {
			break
		}
	}
}

func (c * Client) StopClient() (err error) {
	// todo mettre le cas d'erreur
	close(c.stopChan)
	close(c.mChan)
	return
}