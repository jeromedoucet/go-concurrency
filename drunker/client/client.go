package client

import (
	"encoding/json"
	"go-concurrency/drunker/message"
	"log"
	"strconv"
	"time"
)

// this client is an aggregation of one DBClient (Redis for the moment)
// and one broker client (Nsq here)
// will receive order created from some producer, register it on the Db
// and send it to the waiters through the broker
type Client struct {
	redisCl        DbClient
	brokerProducer BrokerProducer
	topic          string
	stopChan       chan bool
}

type DbClient interface {
	Set(string, interface{}, time.Duration) error
	Get(string) (struct{}, error)
	Remove(string) error
}

type BrokerProducer interface {
	Publish(topic string, body []byte) error
}

// create and start a new client with one DataBase client, one broker client
// the topic to use for the broker and the number of order producer to launch
func StartClient(dbClient DbClient, brokerProducer BrokerProducer, topic string) (c *Client, err error) {
	c = new(Client)
	c.redisCl = dbClient
	c.brokerProducer = brokerProducer
	c.topic = topic
	c.stopChan = make(chan bool, 1)
	go c.listen()
	return
}

func (c *Client) listen() {
	for {
		select {
		case <-c.stopChan:
			log.Println("The client is stopping")
			return
		default:
			o := message.NewOrder(message.NextBeverageType())
			json, _ := json.Marshal(o)
			errR := c.redisCl.Set(strconv.Itoa(int(o.Id)), json, 20)
			if errR != nil {
				log.Printf("error during redis registration: %v", errR)
			} else {
				errB := c.brokerProducer.Publish(c.topic, json)
				if errB != nil {
					log.Printf("error during broker registration: %v", errB)
				}
			}
		}
	}
}

func (c *Client) StopClient() (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovery on some error while trying to close channels : %f", r)
		}
	}()
	if c.stopChan != nil {
		c.stopChan <- true
	}
	return
}
