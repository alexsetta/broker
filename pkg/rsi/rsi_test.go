package rsi

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRSI_Add(t *testing.T) {
	r := NewRSI("teste")
	assert.NotNil(t, r, "The RSI should not be nil")

	r.Add(1.0)
	r.Add(2.0)
	r.Add(3.0)
	r.Add(4.0)
	r.Add(5.0)
	r.Add(6.0)
	r.Add(7.0)
	r.Add(8.0)
	r.Add(9.0)
	r.Add(10.0)
	r.Add(11.0)
	r.Add(12.0)
	r.Add(13.0)
	r.Add(14.0)
	r.Add(15.0)
	r.Add(16.0)

	assert.Equal(t, len(r.prices), 15, "The length of the prices slice should be 15")
	assert.Equal(t, r.prices[0], 2.0, "The first element of the prices slice should be 2.0")
	assert.Equal(t, r.prices[14], 16.0, "The last element of the prices slice should be 16.0")
}

func TestRSI_Calculate(t *testing.T) {
	r := NewRSI("teste")
	assert.NotNil(t, r, "The RSI should not be nil")

	r.Add(31.88)
	r.Add(31.68)
	r.Add(31.81)
	r.Add(32.31)
	r.Add(33.03)
	r.Add(32.91)
	r.Add(32.45)
	r.Add(33.08)
	r.Add(33.27)
	r.Add(33.21)
	r.Add(32.32)
	r.Add(32.97)
	r.Add(32.75)
	r.Add(32.76)
	r.Add(32.59)
	r.Add(31.23)

	fmt.Println(calcRSI(r.prices, 14))

	//rsi := r.Calculate()
	//assert.Equal(t, 19.00, rsi, "The RSI should be 42.66")
}

// Rsi - Relative strength index
func calcRSI(inReal []float64, inTimePeriod int) []float64 {
	outReal := make([]float64, len(inReal))

	if inTimePeriod < 2 {
		return outReal
	}

	// variable declarations
	tempValue1 := 0.0
	tempValue2 := 0.0
	outIdx := inTimePeriod
	today := 0
	prevValue := inReal[today]
	prevGain := 0.0
	prevLoss := 0.0
	today++

	for i := inTimePeriod; i > 0; i-- {
		tempValue1 = inReal[today]
		today++
		tempValue2 = tempValue1 - prevValue
		prevValue = tempValue1
		if tempValue2 < 0 {
			prevLoss -= tempValue2
		} else {
			prevGain += tempValue2
		}
	}
	//wip comparar com a versao original

	prevLoss /= float64(inTimePeriod)
	prevGain /= float64(inTimePeriod)

	tempValue1 = prevGain + prevLoss
	if !((-0.00000000000001 < tempValue1) && (tempValue1 < 0.00000000000001)) {
		outReal[outIdx] = 100.0 * (prevGain / tempValue1)
	} else {
		outReal[outIdx] = 0.0
	}
	outIdx++

	for today < len(inReal) {

		tempValue1 = inReal[today]
		today++
		tempValue2 = tempValue1 - prevValue
		prevValue = tempValue1
		prevLoss *= float64(inTimePeriod - 1)
		prevGain *= float64(inTimePeriod - 1)
		if tempValue2 < 0 {
			prevLoss -= tempValue2
		} else {
			prevGain += tempValue2
		}
		prevLoss /= float64(inTimePeriod)
		prevGain /= float64(inTimePeriod)
		tempValue1 = prevGain + prevLoss
		if !((-0.00000000000001 < tempValue1) && (tempValue1 < 0.00000000000001)) {
			outReal[outIdx] = 100.0 * (prevGain / tempValue1)
		} else {
			outReal[outIdx] = 0.0
		}
		outIdx++
	}

	return outReal
}

func TestRSI_CalculateRSI(t *testing.T) {
	r := NewRSI("ETHBRL")
	assert.NotNil(t, r, "The RSI should not be nil")

	r.Add(6584.92)
	r.Add(6584.92)
	r.Add(6585.64)
	r.Add(6582.45)
	r.Add(6576.54)
	r.Add(6580.43)
	r.Add(6573.94)
	r.Add(6574.50)
	r.Add(6585.09)
	r.Add(6580.40)
	r.Add(6585.78)
	r.Add(6580.00)
	r.Add(6575.01)
	r.Add(6574.43)
	r.Add(6576.82)

	rsi := r.Calculate()
	assert.Equal(t, 42.66, rsi, "The RSI should be 42.66")

	r.Add(6573.53)
	rsi = r.Calculate()
	assert.Equal(t, 40.26, rsi, "The RSI should be 40.26")
}

func TestRSI_CalculateRSIWithFewPrices(t *testing.T) {
	r := NewRSI("ETHBRL")
	assert.NotNil(t, r, "The RSI should not be nil")

	r.Add(6584.92)
	r.Add(6584.92)
	r.Add(6585.64)
	r.Add(6582.45)
	r.Add(6576.54)
	r.Add(6580.43)
	r.Add(6573.94)

	rsi := r.Calculate()
	assert.Equal(t, 0.0, rsi, "The RSI should be 0.0")
}

func TestRSI_ManyRSI(t *testing.T) {
	mr := make(map[string]*RSI)
	mr["ETHBRL"] = NewRSI("ETHBRL")
	assert.NotNil(t, mr["ETHBRL"], "The RSI should not be nil")

	mr["ETHBRL"].Add(6584.92)
	mr["ETHBRL"].Add(6584.92)
	mr["ETHBRL"].Add(6585.64)
	mr["ETHBRL"].Add(6582.45)
	mr["ETHBRL"].Add(6576.54)
	mr["ETHBRL"].Add(6580.43)
	mr["ETHBRL"].Add(6573.94)
	mr["ETHBRL"].Add(6574.50)
	mr["ETHBRL"].Add(6585.09)
	mr["ETHBRL"].Add(6580.40)
	mr["ETHBRL"].Add(6585.78)
	mr["ETHBRL"].Add(6580.00)
	mr["ETHBRL"].Add(6575.01)
	mr["ETHBRL"].Add(6574.43)
	mr["ETHBRL"].Add(6576.82)

	rsi := mr["ETHBRL"].Calculate()
	assert.Equal(t, 42.66, rsi, "The RSI should be 42.66")
}

func TestRSI_LastPrices(t *testing.T) {
	r := NewRSI("BTCUSDT")
	r.LoadPrices()
	assert.NotNil(t, r.prices, "The prices slice should be nil")
	assert.Equal(t, len(r.prices), limit+1, fmt.Sprintf("The length of the prices slice should be %d", limit+1))
	f := r.Calculate()
	assert.NotEqual(t, 0.0, f, "The RSI should be different from 0.0")
	fmt.Println("RSI=", f)
	for _, v := range r.prices {
		fmt.Println(v)
	}
}
