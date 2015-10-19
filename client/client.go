package client
import (
	"time"
	"gopkg.in/redis.v3"
	"xebia.xke.golang/message"
	"xebia.xke.golang/producer"
)

type Client struct {
	mChan chan *message.Order
	stopChan chan struct {}
	redisCl DbClient
}

type DbClient interface{
	Set(string, interface{}, time.Duration) *redis.StatusCmd
	Get(string) *redis.StringCmd
}

func StartClient(dbClient DbClient, countP int) (c *Client, err error) {
	c = new(Client)
	// todo mettre le cas d'erreur
	c.redisCl = dbClient
	c.mChan = make(chan *message.Order)
	c.stopChan = make(chan struct {}, 1)
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
			c.redisCl.Set(string(o.Id), o, 10000000)
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