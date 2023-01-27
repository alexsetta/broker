package rsi

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
)

const (
	limit = 56
	url   = "https://api.binance.us/api/v3/trades?symbol=%s&limit=%d"
)

type RSI struct {
	id     string
	prices []float64
}

type Trade struct {
	Price string `json:"price"`
	Time  int64  `json:"time"`
}

// NewRSI returns a new RSI struct
func NewRSI(id string) *RSI {
	return &RSI{
		id:     id,
		prices: []float64{},
	}
}

// AppendPrice appends a new price to the prices slice
func (r *RSI) Add(price float64) {
	if len(r.prices) == (limit + 1) {
		r.prices = r.prices[1:]
	}
	r.prices = append(r.prices, price)
}

// LoadPrices get last n prices and calculate RSI
func (r *RSI) LoadPrices() {
	r.prices = []float64{}

	// get last n trades
	resp, err := http.Get(fmt.Sprintf(url, r.id, 500))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	// unmarshal json
	var trades []Trade
	err = json.Unmarshal(body, &trades)
	if err != nil {
		return
	}

	// convert value string to float64
	r.prices = []float64{}
	//last := trades[0].Time / 1000
	for _, trade := range trades {
		//secs := trade.Time / 1000
		//diff := int(secs - last)
		//if diff < 5 {
		//	continue
		//}
		//last = secs

		//t := time.Unix(0, trade.Time*int64(time.Millisecond))
		//fmt.Println(t.Format("2006-01-02 15:04:05"), trade.Price, secs)

		f, err := strconv.ParseFloat(trade.Price, 64)
		if err != nil {
			f = 0.0
		}
		r.Add(f)
	}

	//r.prices = r.prices[:15]

	return
}

// Calculate calculates the RSI for the given period
func (r *RSI) Calculate() float64 {
	var avgGain, avgLoss float64

	if len(r.prices) < (limit + 1) {
		return 50.0
	}
	start := len(r.prices) - limit
	finish := len(r.prices)
	interval := finish - start

	for i := start; i < finish; i++ {
		if r.prices[i] > r.prices[i-1] {
			avgGain += r.prices[i] - r.prices[i-1]
		} else {
			avgLoss += r.prices[i-1] - r.prices[i]
		}
	}

	avgGain /= float64(interval)
	avgLoss /= float64(interval)
	rs := avgGain / avgLoss
	rsi := 100 - (100 / (1 + rs))

	if rsi > 100 {
		rsi = 100
	}
	if rsi < 0 {
		rsi = 0
	}

	return math.Round(rsi*100) / 100
}
