package commons

import (
	"fmt"
	"encoding/json"
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

func UnmarshalOrderFromHttp(r *http.Request, order *Order) (buf []byte, err error) {
	buf = make([]byte, r.ContentLength)
	io.ReadFull(r.Body, buf)
	err = json.Unmarshal(buf, &order)
	return
}

//todo unit test
func UnmarshalOrderFromInterface(data interface{}, order *Order) error {
	dec := json.NewDecoder(strings.NewReader(string(data.([]byte))))// todo do thing better
	return dec.Decode(&order)
}
