package main

import (
	"github.com/go-martini/martini"
	"flag"
	"log"
	"encoding/json"
	"gopkg.in/redis.v3"
	"fmt"
	"go-concurrency/messages"
	"go-concurrency/drunker/database"
	"go-concurrency/drunker/client"
)


var redisClient *redis.Client

var dbClient client.DbClient

func main() {
	var redisHost *string = flag.String("redisHost", "127.0.0.1", "redis address ip")
	var redisPort *string = flag.String("redisPort", "6379", "redis port")
	flag.Parse()

	log.Printf("http server listening on port 3000")
	log.Printf("redisHost=%s", *redisHost)
	log.Printf("redisPort=%s", *redisPort)

	connectToRedis(*redisHost, *redisPort)
	initBartenderRestServer()
}


func connectToRedis(redisHost string, redisPort string) {
	var errR error
	dbClient, errR = database.NewRedis(redisHost + ":" + redisPort)
	if errR != nil {
		log.Printf("error during redis connection: %v", errR)
	}
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
		playerId := params["playerId"]
		orderId := params["orderId"]
		fmt.Println(playerId)

		var orderFromRedis message.Order
		orderJsonFromRedis, _ := dbClient.Get(orderId)
		fmt.Println("orderJsonFromRedis=", B2S(orderJsonFromRedis.([]uint8)))
		if orderJsonFromRedis != "" {
			fmt.Println("orderFromRedis!=''")
			//orderFromRedis = message.UmarshallMess(orderJsonFromRedis)
			json.Unmarshal(orderJsonFromRedis.([]byte), &orderFromRedis)
		}

		fmt.Printf("orderFromRedis.Id= %v \n", orderFromRedis)
		if orderFromRedis.PlayerId != "" {
			return 410,  "{ \"status\":\"KO\", \"Error\":\"order has already been requested by :" + orderFromRedis.PlayerId + "\"}"
		}

		orderFromRedis.PlayerId = playerId
		orderAsJson, _ := json.Marshal(orderFromRedis)

		fmt.Println("orderAsJson=", string(orderAsJson))

		err := dbClient.Set(orderId, orderAsJson, 600)
		if err != nil {
			return 500, "{ \"status\":\"KO\", \"Error\":\"" + err.Error() + "\"}"
		}

		return 200, "{ \"status\":\"OK\" }"
	})

	m.Get("/read/:orderId", func(params martini.Params) (int, string) {
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
