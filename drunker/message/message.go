package message

import (
	"fmt"
	"math/rand"
	"time"
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
