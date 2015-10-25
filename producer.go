package main

import (
	"github.com/bitly/go-nsq"
	"go-concurrency/client"
	"log"
	"go-concurrency/database"
)

func main() {
	stp := make(chan *struct{})
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
			<-stp
		}
	}
}
