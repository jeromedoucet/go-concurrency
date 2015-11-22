package main
import (
	"testing"
)

func TestMin(t *testing.T){
	if min(1, 4) != 1 {
		t.Fail()
	}
}


