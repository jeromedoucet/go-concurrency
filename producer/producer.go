package producer

import (
	"xebia.xke.golang/message"
)


type Producer struct {
	stop <-chan struct {}
}


func NewProducer(stop <-chan struct {}) (p * Producer) {
	p = new(Producer)
	p.stop = stop
	return
}

func (p * Producer) Start() <- chan *message.Order {
	c := make(chan *message.Order)
	go func() {
		defer close(c)
		for {
			m := message.NewOrder(message.NextBeverageType())
			select {
			case <-p.stop:
				break
			case c <- m:
			default:
			}
		}
	}()
	return c
}





