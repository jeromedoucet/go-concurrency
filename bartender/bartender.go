package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/vil-coyote-acme/go-concurrency/database"
	"github.com/vil-coyote-acme/go-concurrency/database/redis"
	"github.com/vil-coyote-acme/go-concurrency/messages"
	"log"
	"net/http"
	"runtime"
	"time"
)

var (
	redisHost string
	redisPort string
)

func main() {
	flag.StringVar(&redisHost, "redisHost", "127.0.0.1", "redis address ip")
	flag.StringVar(&redisPort, "redisPort", "6379", "redis port")
	flag.Parse()

	log.Printf("http server listening on port 3000")
	log.Printf("redisHost=%s", redisHost)
	log.Printf("redisPort=%s", redisPort)
	runtime.GOMAXPROCS(runtime.NumCPU())
	// limit the number of idle connections
	http.DefaultTransport.(*http.Transport).MaxIdleConnsPerHost = 100
	initBartenderRestServer()
}

func connectToRedis() database.DbClient {
	dbClient, errR := redis.NewRedis(redisHost + ":" + redisPort)
	if errR != nil {
		log.Panicf("error during redis connection: %v", errR)
	}
	return dbClient
}

func B2S(bs []uint8) string {
	b := make([]byte, len(bs))
	for i, v := range bs {
		if v < 0 {
			b[i] = byte(256 + int(v))
		} else {
			b[i] = byte(v)
		}
	}
	return string(b)
}

// init the bartender Rest server
func initBartenderRestServer() {
	m := martini.Classic()
	//	m.RunOnAddr(":" + port)
	m.Get("/", func(params martini.Params) (int, string) {
		return 200, "hello I'am the bartender"
	})

	m.Post("/bartender/request/:playerId/:orderId", func(params martini.Params) (int, string) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("error : %s", r)
			}
		}()
		playerId := params["playerId"]
		orderId := params["orderId"]
		fmt.Println(playerId)
		dbClient := connectToRedis()
		defer dbClient.Close()

		var orderFromRedis message.Order
		orderJsonFromRedis, err1 := dbClient.Get(orderId)
		if err1 == nil {
			fmt.Println("orderFromRedis!=''")
			json.Unmarshal(orderJsonFromRedis.([]byte), &orderFromRedis)
		} else {
			log.Printf("get an error %s", err1)
			mes := "{ \"status\":\"KO\", \"Error\":\"no order found for " + orderId + "\"}"
			log.Print(mes)
			return 500, mes
		}

		fmt.Printf("orderFromRedis.Id= %v \n", orderFromRedis)
		if orderFromRedis.PlayerId != "" {
			mes := "{ \"status\":\"KO\", \"Error\":\"order has already been requested by :" + orderFromRedis.PlayerId + "\"}"
			log.Print(mes)
			return 410, mes
		}

		orderFromRedis.PlayerId = playerId
		orderAsJson, _ := json.Marshal(orderFromRedis)

		fmt.Println("orderAsJson=", string(orderAsJson))

		err := dbClient.Set(orderId, orderAsJson, 600)
		if err != nil {
			return 500, "{ \"status\":\"KO\", \"Error\":\"" + err.Error() + "\"}"
		}
		wait(orderFromRedis)
		return 200, "{ \"status\":\"OK\" }"
	})

	m.Get("/read/:orderId", func(params martini.Params) (int, string) {
		dbClient := connectToRedis()
		defer dbClient.Close()
		orderId := params["orderId"]
		log.Printf("orderId: %s", orderId)

		jsonMsg, e := dbClient.Get(orderId)
		if e != nil {
			return 500, "{ \"status\":\"KO\", \"Error\":\"" + e.Error() + "\"}"
		}
		fmt.Println(orderId, jsonMsg)

		return 200, "requested order:" + orderId + " = " + B2S(jsonMsg.([]uint8))
	})

	m.Run()
}

func wait(order message.Order) {
	var time2Sleep int64
	switch order.Type {
	case message.Beer:
		time2Sleep = 500
	case message.Cocktail:
		time2Sleep = 10000
	case message.RedWine:
		time2Sleep = 1000
	case message.Vodka:
		time2Sleep = 2500
	case message.Whisky:
		time2Sleep = 2000
	case message.WhiteWine:
		time2Sleep = 1500
	default:
		panic("unrecognized tyope")
	}
	time.Sleep(time.Duration(time.Millisecond * time.Duration(time2Sleep)))
}
