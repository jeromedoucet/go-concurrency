package message
import "math/rand"

type BeverageType int

var cId int = 0
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
	Id int
	PlayerId string
	Type BeverageType
}

func NewOrder(t BeverageType) (o*Order) {
	o=new(Order)
	o.Type = t
	o.Id=cId
	cId ++
	return
}

func NextBeverageType() (t BeverageType) {
	t=BeverageType(rd.Intn(int(Cocktail)))
	return
}