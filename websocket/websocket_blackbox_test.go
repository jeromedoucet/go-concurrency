package websocket_test

import (
	"testing"
	"net/http/httptest"
	"net/http"
	web "github.com/vil-coyote-acme/go-concurrency/websocket"
	"golang.org/x/net/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/vil-coyote-acme/go-concurrency/commons"
	"io/ioutil"
	"encoding/json"
	"strings"
	"time"
)

func Test_websocket_should_push_notification(t *testing.T) {
	// given
	mux := http.NewServeMux()
	notifChan := web.SetupWebsocket(mux)
	srv := httptest.NewServer(mux)
	defer srv.Close()
	// and
	ws, err := websocket.Dial(strings.Replace(srv.URL, "http", "ws", 1) + "/websocket", "", srv.URL)
	if err != nil {
		t.Fatal(err.Error())
	}
	// then
	go func() {
		// when
		notifChan <- commons.Notification{PlayerId:"playerId", Type:commons.Registrate}
		time.Sleep(time.Second * 5)
		ws.Close()
	}()
	var notif commons.Notification
	var data []byte
	data,_ = ioutil.ReadAll(ws) //todo not the right func !
	ws.Close()
	json.Unmarshal(data, &notif)
	assert.Equal(t, "playerId", notif.PlayerId)
	assert.Equal(t, commons.Registrate, notif.Type)
}
