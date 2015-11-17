package message_test

import (
	"go-concurrency/drunker/message"
	"testing"
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
