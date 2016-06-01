package client

import (
	"github.com/vil-coyote-acme/go-concurrency/commons"
	"net/http"
	"time"
)

type RegistrationServer struct {
	RegChan chan commons.RegistrationWrapper
}

func NewRegistrationServer(mux *http.ServeMux) *RegistrationServer {
	res := new(RegistrationServer)
	res.RegChan = make(chan commons.RegistrationWrapper)
	mux.HandleFunc("/registration", res.handleRegistration)
	return res
}

func (rs *RegistrationServer) handleRegistration(w http.ResponseWriter, r *http.Request) {
	var reg commons.Registration
	commons.UnmarshalRegistrationFromHttp(r, &reg)
	regW := commons.RegistrationWrapper{Registration: reg, ResChan: make(chan bool)}
	rs.RegChan <- regW
	// beware here -> the producer must take in account that the consumer may be absent
	res, timedOut := commons.WaitAnswerWithTimeOut(regW.ResChan, time.Second*5)
	if timedOut {
		w.WriteHeader(500)
		return
	}
	if res {
		w.WriteHeader(200)
		return
	}
	w.WriteHeader(403)
}
