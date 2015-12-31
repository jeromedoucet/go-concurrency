package message_test

import (
	"testing"
	"go-concurrency/messages"
)

// basic test
func TestCreateNewOrder(t *testing.T) {
	// quick check of beverage type
	o := message.NewOrder(message.RedWine)
	if o.Type != message.RedWine {
		t.Error("bad type of beverage")
	}
}

// test the beverage type randomization
func TestNextBeverageType(t *testing.T) {
	c := message.NextBeverageType()
	for i := 0; i < 40; i++ {
		c = message.NextBeverageType()
		if message.Beer > c || c > message.Cocktail {
			t.Errorf("The value %d is not in range between %d and %d", c, message.Beer, message.Cocktail)
		}
	}
}

func TestCreateNewOrderCheck(t *testing.T) {
	var id int64 = 12
	player := "la bande a picsou"
	m := message.NewOrderCheck(id, player)
	if m.Id != id {
		t.Errorf("expecting id %d got %d ", id, m.Id)
	} else if m.PlayerId != player {
		t.Errorf("expecting player %s got %s ", player, m.PlayerId)
	}
}

