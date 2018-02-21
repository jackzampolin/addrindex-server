package addrindex

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// NewCurrencyData returns the struct that manages currency data
func NewCurrencyData() *CurrencyData {
	cd := &CurrencyData{}
	cd.Refresh()
	return cd
}

// CurrencyData represents the current BTC price from a couple of providers
type CurrencyData struct {
	Binance        float64 `json:"binance"`
	BlockchainInfo float64 `json:"blockchainInfo"`
	Coinbase       float64 `json:"coinbase"`

	sync.Mutex
}

// JSON returns the json representation of CurrencyData
func (c *CurrencyData) JSON() []byte {
	c.Lock()
	out, _ := json.Marshal(c)
	c.Unlock()
	return out
}

// Refresh refreshes the bitcoin price for the currency info struct
func (c *CurrencyData) Refresh() {
	bn := binancePrice()
	bi := blockchainInfoPrice()
	cb := coinbasePrice()
	c.Lock()
	c.Binance = bn
	c.BlockchainInfo = bi
	c.Coinbase = cb
	c.Unlock()
}

type getBinancePriceResponse struct {
	High      string  `json:"high"`
	Last      string  `json:"last"`
	Timestamp string  `json:"timestamp"`
	Bid       string  `json:"bid"`
	Vwap      string  `json:"vwap"`
	Volume    string  `json:"volume"`
	Low       string  `json:"low"`
	Ask       string  `json:"ask"`
	Open      float64 `json:"open"`
}

type getCoinbasePriceResponse struct {
	Data struct {
		Base     string `json:"base"`
		Currency string `json:"currency"`
		Amount   string `json:"amount"`
	} `json:"data"`
}

func blockchainInfoPrice() float64 {
	req, err := http.NewRequest("GET", "https://blockchain.info/tobtc?currency=usd&value=1000", nil)
	if err != nil {
		log.Println("Failed updating coinbase price - request creation failed")
		return 0.0
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed updating coinbase price - call failed")
		return 0.0
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed updating coinbase price - failed reading body")
		return 0.0
	}

	o, err := strconv.ParseFloat(string(body), 64)
	if err != nil {
		log.Println("Failed updating coinbase price - failed parsing float")
		return 0.0
	}

	return (1 / o) * 1000
}

func coinbasePrice() float64 {
	out := getCoinbasePriceResponse{}
	req, err := http.NewRequest("GET", "https://api.coinbase.com/v2/prices/spot?currency=USD", nil)
	if err != nil {
		log.Println("Failed updating coinbase price - request creation failed")
		return 0.0
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("CB-VERSION", "2015-04-08")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed updating coinbase price - call failed")
		return 0.0
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed updating coinbase price - failed reading body")
		return 0.0
	}

	err = json.Unmarshal(body, &out)
	if err != nil {
		log.Println("Failed updating coinbase price - failed unmarshalling json")
		return 0.0
	}

	o, err := strconv.ParseFloat(out.Data.Amount, 64)
	if err != nil {
		log.Println("Failed updating coinbase price - failed parsing float")
		return 0.0
	}

	return o
}

func binancePrice() float64 {
	out := getBinancePriceResponse{}
	req, err := http.NewRequest("GET", "https://www.bitstamp.net/api/ticker/", nil)
	if err != nil {
		log.Println("Failed updating binance price - request creation failed")
		return 0.0
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Failed updating binance price - call failed")
		return 0.0
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed updating binance price - failed reading body")
		return 0.0
	}

	err = json.Unmarshal(body, &out)
	if err != nil {
		log.Println("Failed updating binance price - failed unmarshalling json")
		return 0.0
	}

	o, err := strconv.ParseFloat(out.Last, 64)
	if err != nil {
		log.Println("Failed updating binance price - failed parsing float")
		return 0.0
	}

	return o
}
