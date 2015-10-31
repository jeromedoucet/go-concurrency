package producer

import (
	"go-concurrency/drunker/message"
	"log"
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
	log.Println("starting a producer")
	go func() {
		for {
			m := message.NewOrder(message.NextBeverageType())
			select {
			case <-p.stop:
			log.Println("stopping producer")
				break
			case c <- m:
			default:
			}
		}
	}()
}





