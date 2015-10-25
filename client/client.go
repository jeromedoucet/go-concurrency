package client
import (
	"time"
	"go-concurrency/message"
	"go-concurrency/producer"
	"encoding/json"
	"log"
	"strconv"
)

// this client is an aggregation of one DBClient (Redis for the moment)
// and one broker client (Nsq here)
// will receive order created from some producer, register it on the Db
// and send it to the waiters through the broker
type Client struct {
	mChan chan *message.Order
	stopChan chan struct {}
	redisCl DbClient
	brokerProducer BrokerProducer
	topic string
}

type DbClient interface{
	Set(string, interface{}, time.Duration) error
	Get(string) (struct{}, error)
}

type BrokerProducer interface {
	Publish(topic string, body []byte) error
}

// create and start a new client with one DataBase client, one broker client
// the topic to use for the broker and the number of order producer to launch
func StartClient(dbClient DbClient, brokerProducer BrokerProducer, topic string, countP int) (c *Client, err error) {
	c = new(Client)
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
			errR :=c.redisCl.Set(strconv.Itoa(int(o.Id)), json, 20)
			if errR != nil {
				log.Printf("error during redis registration: %v",errR)
			} else {
				errB := c.brokerProducer.Publish(c.topic, json)
				if errB != nil {
					log.Printf("error during broker registration: %v",errB)
				}
			}
		}else {
			log.Println("receive nil message. Stop client")
			c.StopClient()
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