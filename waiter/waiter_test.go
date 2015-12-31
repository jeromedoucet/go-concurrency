package main
import (
	"testing"
	"net/http/httptest"
	"net/http"
	"go-concurrency/messages"
	"encoding/json"
	"bytes"
)

func TestCreateDeliverBody(t *testing.T) {
	order := message.NewOrder(message.Beer)
	order.Id = 1234
	b := createDeliverBody(order)
	var body message.OrderCheck
	json.NewDecoder(bytes.NewBuffer(b)).Decode(&body)
	if body.Id != order.Id || body.PlayerId != playerId {
		t.Fail()
	}
}

func TestDeliver(t *testing.T)  {
	c := make(chan bool, 2)
	order := message.NewOrder(message.Beer)
	order.Id = 1234
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() { c <- true }()
		if r.Method != "POST" {
			t.Fail()
		}
		var body message.OrderCheck
		json.NewDecoder(r.Body).Decode(&body)
		if body.Id != order.Id || body.PlayerId != playerId {
			t.Fail()
		}
	}))
	defer ts.Close()
	deliver(ts.URL, createDeliverBody(order))
	select {
	case <-c:
		return
	default:
		t.Fail()
	}
}


func TestDeliveryUrl(t *testing.T) {
	url := deliverUrl("test")
	if url != "http://test/orders" {
		t.Fail()
	}
}

func TestAskBartenderUrl(t *testing.T) {
	order := message.NewOrder(message.Beer)
	order.Id = 1234
	url := askBartenderUrl("test",order)
	if url != "http://test/bartender/request/player/1234" {
		t.Fail()
	}
}

func TestAskBartender(t *testing.T) {

}
