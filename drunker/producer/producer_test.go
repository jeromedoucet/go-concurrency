package producer_test

import (
	"go-concurrency/drunker/message"
	"go-concurrency/drunker/producer"
	"testing"
)

func TestStart(t *testing.T) {
	stp := make(chan struct{}, 1)
	p := producer.NewProducer(stp)
	out := make(chan *message.Order)
	defer close(out)
	defer close(stp)
	p.Start(out)
	order := <-out
	if order == nil {
		t.Errorf("error while testing producer")
	}
	t.Logf("received order : %s", order)
}
