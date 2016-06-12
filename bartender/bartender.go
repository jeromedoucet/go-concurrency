package bartender

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/vil-coyote-acme/go-concurrency/commons"
	"log"
	"net/http"
)

func NewBartender(redisAddr string) *Bartender {
	log.Println(fmt.Sprintf("bartender | create the bartender with the redis addr : %s", redisAddr))
	b := new(Bartender)
	b.redisAddr = redisAddr
	b.mux = http.NewServeMux()
	b.mux.HandleFunc("/orders", b.handleOrder) //mux server. only listen on /order request !
	return b
}

type Bartender struct {
	redisAddr string
	mux       *http.ServeMux
	started   bool
}

func (b *Bartender) Start() {
	if !b.started {
		log.Println("bartender | the bartender is starting, listening on 4343 port")
		b.started = true
		err := http.ListenAndServe(":4343", b.mux)
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}

func (b *Bartender) handleOrder(w http.ResponseWriter, r *http.Request) {
	var order commons.Order
	_, unMarshallErr := commons.UnmarshalOrderFromHttp(r, &order)
	if unMarshallErr != nil {
		log.Printf("An error happends : %s \n\r", unMarshallErr.Error())
		return
	}
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
	_, saveErr := c.Do("SET", existingOrder.Id, string(bd))
	if saveErr != nil {
		log.Printf("An error happends : %s \n\r", saveErr.Error())
		w.WriteHeader(500)
		return
	}
	log.Println(fmt.Sprintf("Bartender | order %s successfully registered", order))
	w.WriteHeader(200)
}
