package producer_test
import (
	"testing"
	"xebia.xke.golang/producer"
)



func TestStart(t *testing.T) {
	stp := make(chan struct{}, 1)
	p := producer.NewProducer(stp)
	out := p.Start();
	order := <-out
	if(order == nil){
		t.Errorf("error while testing producer")
	}
	t.Logf("received order : %s", order)
	close(stp)
}
