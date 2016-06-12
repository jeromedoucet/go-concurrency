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
)

func main() {
	flag.StringVar(&component, "comp", "client", "must be the component name to launch : client or bartender")
	flag.Parse()
	switch component {
	case "bartender":
		bartender.NewBartender(redisAddr).Start()
	case "client":
		client.StartClient(redisAddr)
	default:
		panic(fmt.Sprintf("unknow component. Got : %s", component))
	}
}
