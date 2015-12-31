package main
import (
	"testing"
	"net/http/httptest"
	"net/http"
	"go-concurrency/messages"
	"encoding/json"
	"sync"
)

func TestDeliver(t *testing.T)  {
	var wg sync.WaitGroup
	wg.Add(1)
	order := message.NewOrder(message.Beer)
	order.Id = 1234
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Logf("receive on request %s", r)
		wg.Done()
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
	wg.Wait()
}
