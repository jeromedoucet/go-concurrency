package main
import (
	"log"
	"flag"
	"sync"
	"github.com/nsqio/go-nsq"
	mes "go-concurrency/messages"
	"encoding/json"
	"bytes"
	"net/http"
	"strconv"
)

var (
	lookupaddr string = "51.254.216.243:4161"
	bartenderAddr string = "51.254.216.243:3001"
	deliverAddr string = "51.254.216.243:3002"
	playerId string = "player"
)


func main() {
	topic := flag.String("topic", "orders#ephemeral", "the topic to subscribe on")
	channel := flag.String("channel", "chan#ephemeral", "the channel to use to consume topic message")// to do remove and make it empty
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


func (* Handler) HandleMessage(message *nsq.Message) (e error) {
	defer func() {
		if r := recover(); r != nil {
			e = r.(error)
			return
		}
	}()
	order := unmarshallMes(message)
	resB := askBartender(askBartenderUrl(bartenderAddr, &order))
	if resB == 200 {
		deliver(deliverUrl(deliverAddr), createDeliverBody(&order))
	}
	return
}

func unmarshallMes (message *nsq.Message) mes.Order {
	var order mes.Order
	log.Printf("get the raw message : %s", string(message.Body))
	json.NewDecoder(bytes.NewBuffer(message.Body)).Decode(&order)
	log.Printf("get the order : %s", order)
	return order
}

func askBartender(url string) (statusCode int) {
	resp, err := http.Post(url, "text/plain", bytes.NewBufferString(""))
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
	_, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Panicf("error when trying post on %v ", url)
	}
}