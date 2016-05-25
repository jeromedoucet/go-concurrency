package commons

import "fmt"

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
	Id       int
	Quantity int
	Type     BeverageType
	CallBackUrl string
}

func (o Order) String() string {
	return fmt.Sprintf("id : %d, quantity : %d, type : %d, callback : %s", o.Id, o.Quantity, o.Type, o.CallBackUrl)
}
