package commons

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type BeverageType int

const (
	Beer BeverageType = 0 + iota
	RedWine
	WhiteWine
	Whisky
	Vodka
	Cocktail
)

const (
	Registrate NotificationType = "registrate"
	Unregistrate = "unregistrate"
	Score        = "score"
)

type Order struct {
	Id          int
	Quantity    int
	Type        BeverageType
	CallBackUrl string
	Valid       bool
}

func (o Order) String() string {
	return fmt.Sprintf("id : %d, quantity : %d, type : %d, callback : %s", o.Id, o.Quantity, o.Type, o.CallBackUrl)
}

type Registration struct {
	PlayerId string
	Ip       string
}

func (r Registration) String() string {
	return fmt.Sprintf("PlayerId : %s, Ip : %s", r.PlayerId, r.Ip)
}

type RegistrationWrapper struct {
	Registration
	ResChan chan bool
}

type NotificationType string

type Notification  struct {
	PlayerId string
	Type NotificationType
	Rate float64
	Score int
}

func (n Notification) String() string {
	return fmt.Sprintf("PlayerId : %s, Type : %s, Rate : %d, Score : %d", n.PlayerId, n.Type, n.Rate, n.Score)
}

type Credit struct {
	PlayerId string
	Score int
}

func UnmarshalOrderFromHttp(r *http.Request, order *Order) (buf []byte, err error) {
	buf = make([]byte, r.ContentLength)
	io.ReadFull(r.Body, buf)
	err = json.Unmarshal(buf, &order)
	return
}

//todo unit test
func UnmarshalOrderFromInterface(data interface{}, order *Order) error {
	dec := json.NewDecoder(strings.NewReader(string(data.([]byte)))) // todo do thing better
	return dec.Decode(&order)
}

func UnmarshallCreditFromInterface(data interface{}, credit *Credit) error {
	dec := json.NewDecoder(strings.NewReader(string(data.([]byte)))) // todo do thing better
	return dec.Decode(&credit)
}

//todo unit test
func UnmarshalRegistrationFromHttp(r *http.Request, order *Registration) {
	buf := make([]byte, r.ContentLength)
	io.ReadFull(r.Body, buf)
	json.Unmarshal(buf, &order)
	return
}
