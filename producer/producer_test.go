package producer_test
import (
	"testing"
	"xebia.xke.golang/producer"
	"xebia.xke.golang/message"
)



func TestStart(t *testing.T) {
	stp := make(chan struct{}, 1)
	p := producer.NewProducer(stp)
	out := make(chan *message.Order)
	defer close(out)
	defer close(stp)
	p.Start(out);
	order := <-out
	if(order == nil){
		t.Errorf("error while testing producer")
	}
	t.Logf("received order : %s", order)
}
