package data

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-hclog"
)

func TestNewRates(t *testing.T) {
	tr, err := NewRates(hclog.Default())
	if err != nil {
		t.Fatal(err)
	}
	// fmt.Sprintf("%#v", tr.rates)
	// fmt.Println("=>", tr.rates)
	fmt.Printf("%#v", tr.rates)
}
