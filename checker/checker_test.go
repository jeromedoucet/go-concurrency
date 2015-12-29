package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"
	"go-concurrency/messages"
)

type testReader struct {
	data bytes.Buffer
}

func (s testReader) Read(p []byte) (n int, err error) {
	n, err = s.data.Read(p)
	return
}

func newTestReader(p []byte) testReader {
	t := new(testReader)
	t.data = *bytes.NewBuffer(p)
	return *t
}

type assert func(t *testing.T)

type dbMock struct {
	getKey        string
	t             *testing.T
	badKeyAssert  assert
	goodKeyAssert assert
}

func (m dbMock) Set(key string, value interface{}, ttl time.Duration) error {
	m.t.Error("call to set forbidden !")
	return nil
}

func (m dbMock) Get(key string) (res interface{}, err error) {
	if key == m.getKey {
		m.goodKeyAssert(m.t)
		res = *new(struct{})
	} else {
		m.badKeyAssert(m.t) // call on bad key
		err = *new(error)
	}
	return
}

func (m dbMock) Remove(key string) (err error) {
	return
}

func newDbMock(key string, t *testing.T, badGetAssert assert, goodKeyAssert assert) (d *dbMock) {
	d = new(dbMock)
	d.getKey = key
	d.t = t
	d.badKeyAssert = badGetAssert
	d.goodKeyAssert = goodKeyAssert
	return
}

func TestOnCorrectCheck(t *testing.T) {
	m := newDbMock("12345", t, assertFailure, assertSuccess)
	d := newChecker(*m)
	o := message.NewOrder(message.Beer)
	o.Id = 12345
	json, _ := json.Marshal(o)
	r, _ := http.NewRequest("something", "something", newTestReader(json))
	d.onCheck(*new(http.ResponseWriter), r)
}

func assertFailure(t *testing.T) {
	t.Errorf("get failure !")
}

func TestOnIncorrectCheck(t *testing.T) {
	m := newDbMock("12346", t, assertSuccess, assertFailure)
	d := newChecker(*m)
	o := message.NewOrder(message.Beer)
	o.Id = 12345
	json, _ := json.Marshal(o)
	r, _ := http.NewRequest("something", "something", newTestReader(json))
	d.onCheck(*new(http.ResponseWriter), r)
}

func assertSuccess(t *testing.T) {
}

func TestOnBadRequest(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("the handler should handler any panic")
		}
	}()
	d := newChecker(*new(dbMock))
	r, _ := http.NewRequest("something", "something", newTestReader(make([]byte, 0)))
	d.onCheck(*new(http.ResponseWriter), r)
}
