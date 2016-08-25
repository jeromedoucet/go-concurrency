package websocket

import (
	"golang.org/x/net/websocket"
	"net/http"
	"time"
	"log"
	"fmt"
	"github.com/vil-coyote-acme/go-concurrency/commons"
	"encoding/json"
)

var (
	conChan chan *webSocketWrapper
	wsSlice []*webSocketWrapper
	notifChan chan commons.Notification
)

type webSocketWrapper struct {
	ws    *websocket.Conn
	close chan bool
}

func SetupWebsocket(mux *http.ServeMux) chan commons.Notification {
	wsSlice = make([]*webSocketWrapper, 0)
	conChan = make(chan *webSocketWrapper)
	notifChan = make(chan commons.Notification)
	go handleIncomingConnection()
	mux.Handle("/websocket", websocket.Handler(webSocketHandler))
	return notifChan
}

func handleIncomingConnection() {
	for {
		select {
		case ws := <-conChan:
			wsSlice = append(wsSlice, ws)
		case notif := <-notifChan:
			log.Println(fmt.Sprintf("websocket | handle a new notification %s to forward to everybody", notif))
			data, marshErr := json.Marshal(notif)
			if marshErr == nil {
				publishEvent(data)
			} else {
				log.Println(fmt.Sprintf("websocket | error %s while trying to send notification %s", marshErr.Error(), notif))
			}
		default:
			time.Sleep(time.Second * 5)
		}
	}
}

func publishEvent(data []byte) {
	indexes := make([]int, 0)
	for i, c := range wsSlice {
		_, err := c.ws.Write(data)
		log.Println(fmt.Sprintf("websocket | notification to %s done", c.ws.LocalAddr().String()))
		if err != nil {
			log.Println(fmt.Sprintf("websocket | error %s while trying to write on websocket %s. Closing the connection ", err.Error(), c.ws.RemoteAddr()))
			c.ws.Close()
			c.close <- true
			indexes = append(indexes, i)
		}
	}
	for _, index := range indexes {
		wsSlice = append(wsSlice[:index], wsSlice[index + 1:]...)
	}
}

func webSocketHandler(ws *websocket.Conn) {
	log.Println(fmt.Sprintf("websocket | receive incoming ws connection -> %s", ws.RemoteAddr()))
	wrap := webSocketWrapper{ws:ws, close:make(chan bool)}
	conChan <- &wrap
	<-wrap.close // don't return to avoid closing the connection. Let the client close it!
}

