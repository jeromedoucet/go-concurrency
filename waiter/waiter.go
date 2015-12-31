package main
import (
	"log"
	"flag"
	"sync"
	"github.com/nsqio/go-nsq"
	mes "go-concurrency/messages"
	"encoding/json"
	"bytes"
)

var (
	lookupaddr string = "51.254.216.243:4161"
)


func main() {
	topic := flag.String("topic", "orders#ephemeral", "the topic to subscribe on")
	channel := flag.String("channel", "", "the channel to use to consume topic message")// to do remove and make it empty
	flag.Parse()
	var wg sync.WaitGroup
	wg.Add(1)
	initListener(*topic, *channel)
	wg.Wait()
}

type Handler struct {
}

func initListener(topic, channel string) {
	conf := nsq.NewConfig()
	cons, err := nsq.NewConsumer(topic, channel, conf)
	if err != nil {
		log.Panicf("error when trying to create a consumer for topic : %v and channel : %v", topic, channel)
	}
	// maybe possible to handle message in multiple goroutines
	cons.AddConcurrentHandlers(new(Handler), 5)
	cons.ConnectToNSQLookupd(lookupaddr)
}


func (* Handler) HandleMessage(message *nsq.Message) error {
	var order mes.Order
	log.Printf("get the raw message : %s", string(message.Body))
	json.NewDecoder(bytes.NewBuffer(message.Body)).Decode(&order)
	log.Printf("get the order : %s", order)
	return nil
}