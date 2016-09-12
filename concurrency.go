package main

import (
	"flag"
	"fmt"
	"github.com/vil-coyote-acme/go-concurrency/bartender"
	"github.com/vil-coyote-acme/go-concurrency/client"
)

const (
	redisAddr string = "192.168.99.100:6379"
)

var (
	component string
	myAddr string
	myPort string
)

func main() {
	flag.StringVar(&component, "comp", "client", "must be the component name to launch : client or bartender")
	flag.StringVar(&myAddr, "addr", "http://127.0.0.1", "the 'protocol://ip' value that server must use to contact me !")
	flag.StringVar(&myPort, "port", "4444", "the port exposed by the component")
	flag.Parse()
	switch component {
	case "bartender":
		bartender.NewBartender(redisAddr, myPort).Start()
	case "client":
		client.StartClient(redisAddr, "http://" + myAddr, myPort)
	default:
		panic(fmt.Sprintf("unknow component. Got : %s", component))
	}
}
