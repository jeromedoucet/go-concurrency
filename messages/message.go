package message

import (
	"fmt"
	"math/rand"
	"time"
	"io"
	"bytes"
	"encoding/gob"
	"log"
)

type BeverageType int

var rd = rand.New(rand.NewSource(1))

const (
	Beer BeverageType = 0 + iota
	RedWine
	WhiteWine
	Whisky
	Vodka
	Cocktail
)

type Order struct {
	Id       int64
	PlayerId string
	Valid    bool
	Type     BeverageType
}

type OrderCheck struct {
	Id       int64
	PlayerId string
}

func (o Order) String() string {
	return fmt.Sprintf("order with id %d, playerId %s, and beverage type %d", o.Id, o.PlayerId, o.Type)
}

func NewOrder(t BeverageType) (o *Order) {
	o = new(Order)
	o.Type = t
	o.Id = time.Now().UnixNano()
	return o
}

func NextBeverageType() (t BeverageType) {
	t = BeverageType(rd.Intn(int(Cocktail)))
	return
}

func NewOrderCheck(orderId int64, playerId string) (m *OrderCheck) {
	m = new (OrderCheck)
	m.Id = orderId
	m.PlayerId = playerId
	return
}

func GetReader(t interface{}) (io.Reader) {
	r  := new(bytes.Buffer)
	enc := gob.NewEncoder(r)
	err := enc.Encode(t)
	if err != nil {
		log.Panicf("Error : %f when trying to encode some interface %t into byte array", err, t)
	}
	return r
}
