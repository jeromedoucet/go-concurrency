package main
import (
	"testing"
	"go-concurrency/drunker/client"
)

func TestMin(t *testing.T) {
	if min(1, 4) != 1 {
		t.Fail()
	}
}

func TestGetNbProducer(t *testing.T) {
	clients = make([]*client.Client, 1)
	clients = append(clients, new(client.Client))
	if getNbProducer() != 1 {
		t.Fail()
	}
}

func TestDoRemove(t *testing.T) {
	clients = make([]*client.Client, 1)
	clients = append(clients, new(client.Client), new(client.Client))
	nb := doRemove(1)
	if getNbProducer() != 1 || nb != 1{
		t.Fail()
	}
}

