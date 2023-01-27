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
	limit   = 14
	periods = 14
	url     = "https://api.binance.us/api/v3/trades?symbol=%s&limit=%d"
)

type RSI struct {
	id     string
	prices []float64
}

type Trade struct {
	Price string `json:"price"`
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
	resp, err := http.Get(fmt.Sprintf(url, r.id, limit+1))
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
	for _, trade := range trades {
		f, err := strconv.ParseFloat(trade.Price, 64)
		if err != nil {
			f = 0.0
		}
		r.Add(f)
	}

	return
}

// Calculate calculates the RSI for the given period
func (r *RSI) Calculate() float64 {
	var avgGain, avgLoss float64

	if len(r.prices) < (limit + 1) {
		return 0.0
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

	return math.Round(rsi*100) / 100
}

func (r *RSI) Calculate2() float64 {
	var (
		avgGain, avgLoss float64
		gain, loss       float64
	)

	prices := r.prices
	for i := 1; i < len(prices); i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			gain += change
		} else {
			loss -= change
		}

		if i == periods {
			avgGain = gain / float64(periods)
			avgLoss = loss / float64(periods)
		} else if i > periods {
			avgGain = ((avgGain * float64(periods-1)) + change) / float64(periods)
			avgLoss = ((avgLoss * float64(periods-1)) + change) / float64(periods)
		}
	}

	if avgLoss == 0 {
		return 100
	}

	rs := avgGain / avgLoss
	rsi := 100 - (100 / (1 + rs))

	return rsi
}
