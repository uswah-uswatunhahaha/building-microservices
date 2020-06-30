package data

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strconv"

	"github.com/hashicorp/go-hclog"
)

// ExchangeRates is a Exchange structure with couples elements
type ExchangeRates struct {
	log   hclog.Logger
	rates map[string]float64
}

// NewRates is a constructor to be able get the exchange rates from the European Centre Bank
func NewRates(l hclog.Logger) (*ExchangeRates, error) {
	er := &ExchangeRates{log: l, rates: map[string]float64{}}

	err := er.getRates()

	return er, err
	// return er, nil
}

// GetRate returns a ration between currencies
func (e *ExchangeRates) GetRate(base, dest string) (float64, error) {
	br, ok := e.rates[base]
	if !ok {
		return 0, fmt.Errorf("Rate not found for currency %s", base)
	}
	dr, ok := e.rates[dest]
	if !ok {
		return 0, fmt.Errorf("Rate not found for currency %s", dest)
	}
	return br / dr, nil
}

func (e *ExchangeRates) getRates() error {
	resp, err := http.DefaultClient.Get("https://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
	if err != nil {
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Expected error code 200 got %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	md := &Cubes{}
	xml.NewDecoder(resp.Body).Decode(&md)

	// loop over this collection (Cubes) and then convert them into floats
	// and then put them into our array => rate map[string]float64
	for _, c := range md.CubeData {
		r, err := strconv.ParseFloat(c.Rate, 64)
		if err != nil {
			return err
		}
		e.rates[c.Currency] = r
	}
	// set the euro as a rate
	e.rates["EUR"] = 1

	return nil

}

// Cubes is a collection of Cube
type Cubes struct {
	CubeData []Cube `xml:"Cube>Cube>Cube"`
}

// Cube have two attributes: Currency and Rate
type Cube struct {
	Currency string `xml:"currency,attr"`
	Rate     string `xml:"rate,attr"`
}
