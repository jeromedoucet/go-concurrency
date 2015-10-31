package main

import (
	"github.com/bitly/go-nsq"
	"go-concurrency/drunker/client"
	"log"
	"go-concurrency/drunker/database"
	"flag"
)


func main() {
	nbProducer := flag.Int("nbProducer", 1, "number of producer to run")
	flag.Parse()
	log.Printf("GO-CONCURRENCY producer module is starting with %d prducer", *nbProducer)
	stp := make(chan *struct{})
	for i:= 0; i < *nbProducer; i ++ {
		startOneProducer()
	}
	<-stp

}

func startOneProducer() {
	config := nsq.NewConfig()
	w, errN := nsq.NewProducer("127.0.0.1:4150", config)
	if errN != nil {
		log.Printf("error during nsq producer creation: %v",errN)
	} else {
		d, errR := database.NewRedis("127.0.0.1:6379")
		if errR != nil {
			log.Printf("error during redis connection: %v",errR)
		} else {
			client.StartClient(d, w, "orders#ephemeral", 1)
		}
	}
}
