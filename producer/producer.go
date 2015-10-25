package producer

import (
	"go-concurrency/message"
)


type Producer struct {
	stop <-chan struct {}
}


func NewProducer(stop <-chan struct {}) (p *Producer) {
	p = new(Producer)
	p.stop = stop
	return
}

func (p * Producer) Start(c chan *message.Order) {
	go func() {
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
}





