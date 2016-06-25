package client

import (
	"net/http"
	"github.com/vil-coyote-acme/go-concurrency/commons"
	"strings"
	"github.com/garyburd/redigo/redis"
	"log"
	"math"
	"encoding/json"
	"fmt"
)

var (
	redisAddr string
	ch chan commons.Notification
)

func initPaymentHandling(mux *http.ServeMux, redis string, notifC chan commons.Notification) {
	mux.HandleFunc("/bill/", paymentEndPoint)
	redisAddr = redis
	ch = notifC
}

func paymentEndPoint(w http.ResponseWriter, r *http.Request) {
	var order commons.Order
	var credit commons.Credit
	urlPart := strings.Split(r.URL.Path, "/")
	// todo test me
	if len(urlPart) < 3 {
		log.Printf("payment | An error happends. Not enought variable path. Expected 3, get %d \n\r", len(urlPart))
		w.WriteHeader(404)
		return
	}
	c, redisErr := redis.Dial("tcp", redisAddr)
	defer c.Close()
	// todo test me
	if redisErr != nil {
		log.Printf("payment | An error happends : %s \n\r", redisErr.Error())
		w.WriteHeader(500)
		return
	}
	// check the existance of the order
	exist, existErr := c.Do("EXISTS", urlPart[3])
	if existErr != nil {
		log.Printf("An error happends : %s \n\r", existErr.Error())
		w.WriteHeader(500)
		return
	}
	// todo test me
	if exist.(int64) != 1 {
		log.Printf("Any order founded for order %d \n\r", urlPart[3])
		w.WriteHeader(404)
		return
	}
	data, getError := c.Do("GET", urlPart[3])
	if getError != nil {
		log.Printf("An error happends : %s \n\r", getError.Error())
		w.WriteHeader(500)
		return
	}
	// check that the order is not already validated
	// todo test me
	commons.UnmarshalOrderFromInterface(data, &order)
	if !order.Valid {
		// not validated
		log.Printf("Order with id %d is not validated \n\r", order.Id)
		w.WriteHeader(403)
		return
	}
	c.Do("DEL", order.Id)
	// check the existance of the order
	existCredit, _ := c.Do("EXISTS", urlPart[2])
	// todo test me
	if existCredit.(int64) != 1 {
		credit = commons.Credit{PlayerId:urlPart[2], Score:computePaymentSum(order)}
	} else {
		dataCred, _ := c.Do("GET", urlPart[2])
		commons.UnmarshallCreditFromInterface(dataCred, &credit)
		credit.Score += computePaymentSum(order)
	}

	bd, _ := json.Marshal(credit)
	_, saveErr := c.Do("SET", credit.PlayerId, string(bd))
	if saveErr != nil {
		log.Printf("An error happends : %s \n\r", saveErr.Error())
		w.WriteHeader(500)
		return
	}
	log.Println(fmt.Sprintf("payment | credit %s successfully registered", credit))
	func () {ch <- commons.Notification{PlayerId:credit.PlayerId, Type:commons.Score, Score:credit.Score}} ()
	w.WriteHeader(200)
}


func computePaymentSum(order commons.Order) int {
	return int(math.Pow(float64(order.Type + 1), 2.0)) * order.Quantity
}
