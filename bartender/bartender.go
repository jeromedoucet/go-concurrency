package bartender

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/vil-coyote-acme/go-concurrency/commons"
	"log"
	"net/http"
	"time"
	"math"
)

var (
	httpPort string
)

func NewBartender(redisAddr string, myPort string) *Bartender {
	log.Println(fmt.Sprintf("bartender | create the bartender with the redis addr : %s", redisAddr))
	httpPort = myPort
	b := new(Bartender)
	b.redisAddr = redisAddr
	b.mux = http.NewServeMux()
	b.mux.HandleFunc("/orders", b.handleOrder) //mux server. only listen on /order request !
	b.tokenBuf = make(map[string]chan bool)
	b.tokenChan = make(chan tokenReq)
	return b
}

type Bartender struct {
	redisAddr string
	mux       *http.ServeMux
	started   bool
	tokenBuf  map[string]chan bool //use to avoid too many conection from 1 server
	tokenChan chan tokenReq
}

type tokenReq struct {
	playerId string
	res      chan chan bool
}

func (b *Bartender) Start() {
	if !b.started {
		log.Println(fmt.Sprintf("bartender | the bartender is starting, listening on %s port", httpPort))
		b.started = true
		b.tokenProviderLoop()
		err := http.ListenAndServe(":" + httpPort, b.mux)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

func (b *Bartender) tokenProviderLoop() {
	go func() {
		for {
			req := <-b.tokenChan
			c, p := b.tokenBuf[req.playerId]
			if !p {
				c = make(chan bool, 5)
				for i := 0; i < 5; i ++ {
					c <- true
				}
				b.tokenBuf[req.playerId] = c
			}
			req.res <- c
		}
	}()
}

func (b *Bartender) handleOrder(w http.ResponseWriter, r *http.Request) {
	var order commons.Order
	unMarshallErr := commons.UnmarshalOrderFromHttp(r, &order)
	if unMarshallErr != nil {
		log.Printf("An error happends : %s \n\r", unMarshallErr.Error())
		return
	}
	// try to get one token for the
	req := tokenReq{playerId:order.PlayerId, res:make(chan chan bool)}
	b.tokenChan <- req
	c := <-req.res
	select {
	case token := <-c:
	// one token is available
		defer func() {
			c <- token
		}()
		b.doHandleOrder(w, r, order)
	default:
		b.doApplyPenalty(w, r, order)
	}

}

func (b*Bartender) doApplyPenalty(w http.ResponseWriter, r *http.Request, order commons.Order) {
	var credit commons.Credit
	c, redisErr := redis.Dial("tcp", b.redisAddr)
	defer c.Close()
	// todo test me
	if redisErr != nil {
		log.Printf("payment | An error happends : %s \n\r", redisErr.Error())
		w.WriteHeader(500)
		return
	}
	existCredit, _ := c.Do("EXISTS", order.PlayerId)
	// todo test me
	if existCredit.(int64) != 1 {
		credit = commons.Credit{PlayerId:order.PlayerId, Score: -10000, Timestamp:int(time.Now().UTC().Unix())}
	} else {
		dataCred, _ := c.Do("GET", order.PlayerId)
		commons.UnmarshallCreditFromInterface(dataCred, &credit)
		credit.Score += -10000
	}

	bd, _ := json.Marshal(credit)
	_, saveErr := c.Do("SET", credit.PlayerId, string(bd))
	if saveErr != nil {
		log.Printf("An error happends : %s \n\r", saveErr.Error())
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(403)
}

func (b*Bartender) doHandleOrder(w http.ResponseWriter, r *http.Request, order commons.Order) {
	log.Println(fmt.Sprintf("Bartender | receive one order : %s", order))
	c, redisErr := redis.Dial("tcp", b.redisAddr)
	defer c.Close()
	if redisErr != nil {
		log.Printf("An error happends : %s \n\r", redisErr.Error())
		w.WriteHeader(500)
		return
	}
	// check the existance of the order
	exist, existErr := c.Do("EXISTS", order.Id)
	if existErr != nil {
		log.Printf("An error happends : %s \n\r", existErr.Error())
		w.WriteHeader(500)
		return
	}
	if exist.(int64) != 1 {
		log.Printf("Any order founded for order %d \n\r", order.Id)
		w.WriteHeader(404)
		return
	}
	// get the existing order
	var existingOrder commons.Order
	data, getError := c.Do("GET", order.Id)
	if getError != nil {
		log.Printf("An error happends : %s \n\r", getError.Error())
		w.WriteHeader(500)
		return
	}
	// check that the order is not already validated
	commons.UnmarshalOrderFromInterface(data, &existingOrder)
	if existingOrder.Valid {
		// already validated
		log.Printf("Order with id %d is already validated \n\r", order.Id)
		w.WriteHeader(403)
		return
	}
	// register the validated
	existingOrder.Valid = true
	bd, marshalError := json.Marshal(existingOrder)
	if marshalError != nil {
		log.Printf("An error happends : %s \n\r", marshalError.Error())
		w.WriteHeader(500)
		return
	}
	d := time.Duration(int(math.Pow(float64(order.Type + 1), 2.0)) * order.Quantity) * time.Millisecond * 10
	log.Println(fmt.Sprintf("Bartender | wait for %d millisecond", d))
	time.Sleep(d)
	_, saveErr := c.Do("SET", existingOrder.Id, string(bd))
	if saveErr != nil {
		log.Printf("An error happends : %s \n\r", saveErr.Error())
		w.WriteHeader(500)
		return
	}
	log.Println(fmt.Sprintf("Bartender | order %s successfully registered", order))
	w.WriteHeader(200)
}
