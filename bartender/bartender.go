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
	"io"
)


var redisClient *redis.Client

var dbClient client.DbClient

//type SomeData struct {
//	// see http://mholt.github.io/json-to-go/
//	ID string `json:"id"`
//	Name string `json:"name"`
//}

//type Order struct {
//	Id int64 `json:"id"`
//	PlayerId string `json:"playerId"`
//	Valid bool `json:"valid"`
//	Type string `json:"type"`
//}


func main() {
	var httpServerPort *string = flag.String("http-port", "3000", "rest interface listening port")
//	var redisAddr *string = flag.String("redis", "127.0.0.1:6379", "redis address ip:port")
	var redisHost *string = flag.String("redisHost", "127.0.0.1", "redis address ip")
	var redisPort *string = flag.String("redisPort", "6379", "redis port")
	flag.Parse()

	log.Printf("Http server configured on port %s", *httpServerPort)
//	log.Printf("Redis server should be on address %s", *redisAddr)

//	initRedisClient(*redisAddr)
	connectToRedis(*redisHost, *redisPort)
	initBartenderRestServer(*httpServerPort)
}

func connectToRedis(redisHost string, redisPort string) {
	var errR error
	dbClient, errR = database.NewRedis(redisHost + ":" + redisPort)
	if errR != nil {
		log.Printf("error during redis connection: %v", errR)
	}
}

//func initRedisClient(redisAddress string) {
//	fmt.Println("redisAddress received: %s", redisAddress)
//	redisClient = redis.NewClient(&redis.Options{
//	Addr:     "localhost:6379",
//		Password: "", // no password set
//		DB:       0,  // use default DB
//	})
//
//	pong, err := redisClient.Ping().Result()
//	fmt.Println(pong, err)
//}

func umarshallMess(from io.Reader, to interface{}) {
	err := json.NewDecoder(from).Decode(to)
	if err != nil {
		log.Panicf("Error when trying to decode request body : %s", err.Error())
	}
}

// init the bartender Rest server
func initBartenderRestServer(port string) {
	m := martini.Classic()
	//	m.RunOnAddr(":" + port)


	m.Get("/", func(params martini.Params) (int, string) {
		return 200, "hello I'am the bartender"
	})


//	// uniquement pour les tests, les orders devraient être insérés dans redis par le drunker/producer
//	m.Post("/bartender/write/order", binding.Bind(Order{}), func(params martini.Params, order Order) (int, string) {
//		log.Println("order=", order)
//
//		log.Println("order.Id:",order.Id)
//
//		orderAsJson, _ := json.Marshal(order)
//
//		err := redisClient.Set(strconv.Itoa(int(order.Id)), orderAsJson, 0).Err()
//		if err != nil {
//			panic(err)
//		}
//
//		return 200, "order:" + strconv.Itoa(int(order.Id)) + " written in redis:" + string(orderAsJson)
//	})


	m.Post("/bartender/request/:playerId/:orderId", func(params martini.Params) (int, string) {
		playerId := params["playerId"]
		orderId := params["orderId"]
		fmt.Println(playerId)
//		orderIdExists, errRead := redisClient.Exists(orderId).Result()
//		if !orderIdExists {
//			return 404, "{ \"status\":\"KO\", \"Error\":\"No order with id " + orderId + "\"}"
//		}

//		orderJsonFromRedis, errRead := redisClient.Get(orderId).Result()
		orderJsonFromRedis, _ := redisClient.Get(orderId).Result()
//		fmt.Println("orderjson=", orderJsonFromRedis)
//		if errRead != nil { // TODO: comment connaitre le type d'erreur ?
//			return 500, "{ \"status\":\"KO\", \"Error\":\"" + errRead.Error() + "\"}"
//		}

		fmt.Println("orderJsonFromRedis=", orderJsonFromRedis)
		var orderFromRedis message.Order
		if orderJsonFromRedis != "" {
			fmt.Println("orderFromRedis!=''")
			//orderFromRedis = message.UmarshallMess(orderJsonFromRedis)
			json.Unmarshal([]byte(orderJsonFromRedis), &orderFromRedis)
		}

		fmt.Printf("orderFromRedis.Id= %v \n", orderFromRedis)
		if orderFromRedis.PlayerId != "" {
			return 410,  "{ \"status\":\"KO\", \"Error\":\"order has already been requested by :" + orderFromRedis.PlayerId + "\"}"
		}

		orderFromRedis.PlayerId = playerId
		orderAsJson, _ := json.Marshal(orderFromRedis)

		fmt.Println("orderAsJson=", string(orderAsJson))

		err := redisClient.Set(orderId, orderAsJson, 0).Err()
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
//			log.Panic(e)
			return 500, "{ \"status\":\"KO\", \"Error\":\"" + e.Error() + "\"}"
		}
		fmt.Println(orderId, jsonMsg)

		return 200, "requested order:" + orderId + " = " //  + jsonMsg
	})



	m.Run()
}
