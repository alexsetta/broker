package rsi

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	dirFile = "d:/dev/go/app/coletor/files"
)

func TestRSI_Add(t *testing.T) {
	r := NewRSI("teste", "", false)
	assert.Equal(t, &RSI{id: "teste", filePath: "", loadFile: false}, r)

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

func TestRSI_CalculateRSI(t *testing.T) {
	r := NewRSI("ETHBRL", dirFile, false)
	assert.Equal(t, &RSI{id: "ETHBRL", filePath: dirFile, loadFile: false}, r)

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
	r := NewRSI("ETHBRL", dirFile, false)
	assert.Equal(t, &RSI{id: "ETHBRL", filePath: dirFile, loadFile: false}, r)

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
	mr["ETHBRL"] = NewRSI("ETHBRL", dirFile, false)
	assert.Equal(t, &RSI{id: "ETHBRL", filePath: dirFile, loadFile: false}, mr["ETHBRL"])

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

func TestRSI_LoadRSI(t *testing.T) {
	// Antes é necessário rodar "TestRSI_CalculateRSI" para gerar o arquivo
	TestRSI_CalculateRSI(t)

	r := NewRSI("ETHBRL", dirFile, true)
	assert.Equal(t, &RSI{id: "ETHBRL", filePath: dirFile, loadFile: true}, r)

	rsi := r.Calculate()
	assert.Equal(t, 40.26, rsi, "The RSI should be 40.26")
}
