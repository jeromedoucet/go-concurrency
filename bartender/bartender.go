package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"gopkg.in/redis.v3"
	"log"
)

var client *redis.Client

type SomeData struct {
	// see http://mholt.github.io/json-to-go/
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	var httpServerPort *string = flag.String("http-port", "3000", "rest interface listening port")
	var redisAddr *string = flag.String("redis", "127.0.0.1:6379", "redis address ip:port")
	flag.Parse()

	log.Printf("Http server configured on port %s", *httpServerPort)
	log.Printf("Redis server should be on address %s", *redisAddr)

	initRedisClient(*redisAddr)
	initBartenderRestServer(*httpServerPort)
}

func initRedisClient(redisAddress string) {
	fmt.Println("am I here ?")
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
}

// init the bartender Rest server
func initBartenderRestServer(port string) {
	m := martini.Classic()
	//	m.RunOnAddr(":" + port)

	m.Post("/bartender/:orderId", binding.Bind(SomeData{}), func(params martini.Params, data SomeData) (int, string) {
		orderId := params["orderId"]
		log.Printf("orderId: %s", orderId)
		log.Println("data=", data)

		log.Println("data.name:", data.Name)

		jsonData, _ := json.Marshal(data)

		err := client.Set(orderId, jsonData, 0).Err()
		if err != nil {
			panic(err)
		}

		return 200, "requested order:" + orderId + " setting data=" + string(jsonData)
	})

	m.Get("/", func(params martini.Params) (int, string) {
		return 200, "hello I'am the bartender"
	})

	m.Get("/read/:orderId", func(params martini.Params) (int, string) {
		orderId := params["orderId"]
		log.Printf("orderId: %s", orderId)

		val, err := client.Get(orderId).Result()
		if err != nil {
			panic(err)
		}
		fmt.Println(orderId, val)

		return 200, "requested order:" + orderId + " = " + val
	})

	m.Run()
}
