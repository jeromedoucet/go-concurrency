package client

import (
	"github.com/vil-coyote-acme/go-concurrency/commons"
	"net/http"
	"time"
	"log"
	"fmt"
	"github.com/vil-coyote-acme/go-concurrency/websocket"
	"os"
)

var (
	regChan      chan commons.RegistrationWrapper
	registration map[string]commons.Registration
	notifChan chan commons.Notification
	started bool
	MyAddr string
	MyPort string
)

func StartClient(redisAddr string, myAddr string, myPort string) {
	log.Println(fmt.Sprintf("client | create the client with the redis addr : %s and the client addr : %s", redisAddr, myAddr))
	MyAddr = myAddr
	MyPort = myPort
	initRegistration(redisAddr)
	mux := http.NewServeMux()
	initRegistrationHandling(mux)
	initWebServer(mux)
	notifChan = websocket.SetupWebsocket(mux)
	initPaymentHandling(mux, redisAddr, notifChan)
	if !started {
		log.Println("client | the client is starting, listening on 4444 port")
		started = true
		err := http.ListenAndServe(":" + myPort, mux)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

}

func initRegistration(redisAddr string) {
	regChan = make(chan commons.RegistrationWrapper)
	registration = make(map[string]commons.Registration, 10)
	go handleRegistration(redisAddr)
}

func initRegistrationHandling(mux *http.ServeMux) {
	mux.HandleFunc("/registration", registrationEndPoint)
}

func initWebServer(mux *http.ServeMux) {
	fs := http.FileServer(http.Dir(os.Getenv("GOPATH") + "/src/github.com/vil-coyote-acme/go-concurrency/web/"))
	mux.Handle("/", fs)
}

func registrationEndPoint(w http.ResponseWriter, r *http.Request) {
	var reg commons.Registration
	commons.UnmarshalRegistrationFromHttp(r, &reg)
	log.Println(fmt.Sprintf("client | receive the registration query : %s", reg))
	regW := commons.RegistrationWrapper{Registration: reg, ResChan: make(chan bool)}
	regChan <- regW
	// beware here -> the producer must take in account that the consumer may be absent
	res, timedOut := commons.WaitAnswerWithTimeOut(regW.ResChan, time.Second * 5)
	if timedOut {
		log.Println(fmt.Sprintf("client | registration unknow issue on query : %s", reg))
		w.WriteHeader(500)
		return
	}
	if res {
		log.Println(fmt.Sprintf("client | registration successful on query : %s", reg))
		w.WriteHeader(200)
		return
	}
	log.Println(fmt.Sprintf("client | registration failed for consistency issue on query : %s", reg))
	w.WriteHeader(403)
}

func handleRegistration(redisAddr string) {
	unRegChan := make(chan commons.Registration)
	for {
		select {
		case rw := <-regChan:
			noConflict := hasNoConflict(&rw.Registration)
			if noConflict {
				registration[rw.PlayerId] = rw.Registration
				notifChan <- commons.Notification{PlayerId:rw.PlayerId, Type:commons.Registrate}
				startNewOrderMaker(MyAddr + ":" + MyPort, redisAddr, rw.Registration, unRegChan)
			}
			rw.ResChan <- noConflict
		case unReg := <-unRegChan:
			delete(registration, unReg.PlayerId)
			notifChan <- commons.Notification{PlayerId:unReg.PlayerId, Type:commons.Unregistrate}
		default:
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func hasNoConflict(r *commons.Registration) (res bool) {
	res = true
	re, ex := registration[r.PlayerId]
	if ex {
		res = re.Ip == r.Ip
		return
	}
	for _, val := range registration {
		if val.Ip == r.Ip {
			res = false
			break
		}
	}
	return
}
