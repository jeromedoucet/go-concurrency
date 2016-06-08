package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/vil-coyote-acme/go-concurrency/commons"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
)

func Test_registration_handling_should_create_produce_registration_wrapper_msg(t *testing.T) {
	// given
	r := commons.Registration{Ip: "http://my-addr:1234", PlayerId: "id"}
	body, _ := json.Marshal(r)
	mux := http.NewServeMux()
	regChan = make(chan commons.RegistrationWrapper)
	initRegistrationHandling(mux)
	s := httptest.NewServer(mux)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	defer s.Close()
	// when
	go func() {
		res, reqErr := http.Post(s.URL + "/registration", "application/json", bytes.NewBuffer(body))
		assert.Nil(t, reqErr)
		assert.Equal(t, 200, res.StatusCode)
		wg.Done()
	}()
	wr, err := waitRegistrationTimeout(regChan, time.Second * 5)
	assert.Nil(t, err)
	assert.Equal(t, r.Ip, wr.Ip)
	assert.Equal(t, r.PlayerId, wr.PlayerId)
	wr.ResChan <- true
	timedOut := commons.WaitTimeout(wg, time.Second * 5)
	assert.False(t, timedOut)
}

func Test_registration_handling_should_return_500_when_no_answer_from_registration_core(t *testing.T) {
	// given
	r := commons.Registration{Ip: "http://my-addr:1234", PlayerId: "id"}
	body, _ := json.Marshal(r)
	mux := http.NewServeMux()
	regChan = make(chan commons.RegistrationWrapper)
	initRegistrationHandling(mux)
	s := httptest.NewServer(mux)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	defer s.Close()
	// when
	go func() {
		res, reqErr := http.Post(s.URL + "/registration", "application/json", bytes.NewBuffer(body))
		assert.Nil(t, reqErr)
		assert.Equal(t, 500, res.StatusCode)
		wg.Done()
	}()
	wr, err := waitRegistrationTimeout(regChan, time.Second * 5)
	assert.Nil(t, err)
	assert.Equal(t, r.Ip, wr.Ip)
	assert.Equal(t, r.PlayerId, wr.PlayerId)
	timedOut := commons.WaitTimeout(wg, time.Second * 10)
	assert.False(t, timedOut)
}

func Test_registration_handling_should_return_403_when_refusal_from_registration_core(t *testing.T) {
	// given
	r := commons.Registration{Ip: "http://my-addr:1234", PlayerId: "id"}
	body, _ := json.Marshal(r)
	mux := http.NewServeMux()
	regChan = make(chan commons.RegistrationWrapper)
	initRegistrationHandling(mux)
	s := httptest.NewServer(mux)
	wg := new(sync.WaitGroup)
	wg.Add(1)
	defer s.Close()
	// when
	go func() {
		res, reqErr := http.Post(s.URL + "/registration", "application/json", bytes.NewBuffer(body))
		assert.Nil(t, reqErr)
		assert.Equal(t, 403, res.StatusCode)
		wg.Done()
	}()
	wr, err := waitRegistrationTimeout(regChan, time.Second * 5)
	assert.Nil(t, err)
	assert.Equal(t, r.Ip, wr.Ip)
	assert.Equal(t, r.PlayerId, wr.PlayerId)
	wr.ResChan <- false
	timedOut := commons.WaitTimeout(wg, time.Second * 5)
	assert.False(t, timedOut)
}

func Test_registration_core_should_accept_new_registration(t *testing.T) {
	// given
	r := commons.Registration{Ip: "http://my-addr:1234", PlayerId: "id"}
	resChan := make(chan bool)
	rw := commons.RegistrationWrapper{Registration:r, ResChan: resChan}
	InitRegistration()
	wg := new(sync.WaitGroup)
	wg.Add(1)
	// when
	go func() {
		regChan <- rw
		wg.Done()
	}()
	timedOut1 := commons.WaitTimeout(wg, time.Second * 5)
	assert.False(t, timedOut1)
	res, timedOut2 := commons.WaitAnswerWithTimeOut(resChan, time.Second * 5)
	assert.False(t, timedOut2)
	assert.True(t, res)
}

func Test_registration_core_should_refuse_acception_when_ip_already_registered(t *testing.T) {

}

func Test_registration_core_should_refuse_acception_when_playerId_already_registered(t *testing.T) {

}

func waitRegistrationTimeout(c chan commons.RegistrationWrapper, timeout time.Duration) (res commons.RegistrationWrapper, err error) {
	select {
	case res = <-c:
		return
	case <-time.After(timeout):
		err = errors.New("time out !")
		return
	}
}
