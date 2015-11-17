package main

import (
	"flag"
	"github.com/bitly/go-nsq"
	"go-concurrency/drunker/client"
	"go-concurrency/drunker/database"
	"log"
)

var (
	nbProducer int
	nsqHost    string
	nsqPort    string
	redisHost  string
	redisPort  string
)

func main() {
	flag.IntVar(&nbProducer, "nbProducer", 1, "number of producer to run")
	flag.StringVar(&nsqHost, "nsqHost", "127.0.0.1", "nsq host")
	flag.StringVar(&nsqPort, "nsqPort", "4150", "nsq port")
	flag.StringVar(&redisHost, "redisHost", "127.0.0.1", "redis host")
	flag.StringVar(&redisPort, "redisPort", "6379", "redis port")
	flag.Parse()
	log.Printf("GO-CONCURRENCY producer module is starting with %d prducer", nbProducer)
	stp := make(chan *struct{})
	for i := 0; i < nbProducer; i++ {
		startOneProducer()
	}
	<-stp

}

func startOneProducer() {
	config := nsq.NewConfig()
	w, errN := nsq.NewProducer(nsqHost+":"+nsqPort, config)
	if errN != nil {
		log.Printf("error during nsq producer creation: %v", errN)
	} else {
		d, errR := database.NewRedis(redisHost + ":" + redisPort)
		if errR != nil {
			log.Printf("error during redis connection: %v", errR)
		} else {
			client.StartClient(d, w, "orders#ephemeral", 1)
		}
	}
}
