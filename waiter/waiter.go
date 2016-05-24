package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/nsqio/go-nsq"
	mes "github.com/vil-coyote-acme/go-concurrency/messages"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var (
	lookupaddr    string           = "127.0.0.1:4161"
	bartenderAddr string           = "127.0.0.1:3000"
	deliverAddr   string           = "127.0.0.1:3002"
	messageChan   chan nsq.Message = make(chan nsq.Message, 20)
	playerId      string
	topic         string
)

func main() {
	flag.StringVar(&playerId, "player", "foo", "the user name")
	flag.StringVar(&topic, "topic", "orders", "the topic to subscribe on")
	flag.Parse()
	channel := playerId
	var wg sync.WaitGroup
	wg.Add(1)
	initListener(topic, channel)
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
	for i := 0; i < 20; i++ {
		go asyncTreatment()
	}
	// maybe possible to handle message in multiple goroutines
	cons.AddConcurrentHandlers(new(Handler), 5)
	cons.ConnectToNSQLookupd(lookupaddr)
}

func (*Handler) HandleMessage(message *nsq.Message) (e error) {
	messageChan <- *message
	return
}

func asyncTreatment() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("receive one error %s", r)
			asyncTreatment()
		}
	}()
	for {
		message := <-messageChan
		log.Printf("receive a message %v", message)
		order := unmarshallMes(&message)
		resB := askBartender(askBartenderUrl(bartenderAddr, &order))
		log.Printf("receive a response from bartender %d", resB)
		if resB == 200 {
			deliver(deliverUrl(deliverAddr), createDeliverBody(&order))
		}
	}
}

func unmarshallMes(message *nsq.Message) mes.Order {
	var order mes.Order
	log.Printf("get the raw message : %s", string(message.Body))
	json.NewDecoder(bytes.NewBuffer(message.Body)).Decode(&order)
	log.Printf("get the order : %s", order)
	return order
}

func askBartender(url string) (statusCode int) {
	resp, err := http.Post(url, "text/plain", bytes.NewBufferString(""))
	resp.Body.Close()
	if err != nil {
		log.Panicf("error when trying to send post on %v ", url)
	} else {
		statusCode = resp.StatusCode
	}
	return
}

func askBartenderUrl(host string, order *mes.Order) string {
	return "http://" + host + "/bartender/request/" + playerId + "/" + strconv.Itoa(int(order.Id))
}

func createDeliverBody(order *mes.Order) []byte {
	o := mes.NewOrderCheck(order.Id, playerId)
	b, err := json.Marshal(o)
	if err != nil {
		log.Panicf("error when trying serialise %s ", order)
	}
	return b
}

func deliverUrl(host string) string {
	return "http://" + host + "/orders"
}

func deliver(url string, body []byte) {
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	resp.Body.Close()
	if err != nil {
		log.Panicf("error when trying post on %v ", url)
	}
}
