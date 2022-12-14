package util

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestUSDToBRL(t *testing.T) {
	s, err := GetHttp("https://economia.awesomeapi.com.br/json/last/usd-brl")
	if err != nil {
		panic(err)
	}

	type result struct {
		USDBRL struct {
			Bid string `json:"bid"`
		} `json:"USDBRL"`
	}
	var a result
	err = json.Unmarshal([]byte(s), &a)
	if err != nil {
		panic(err)
	}

	fmt.Println(a)
	fmt.Println(StringToValue(fmt.Sprintf("%v", a.USDBRL.Bid)))
}
