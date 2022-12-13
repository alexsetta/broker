package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/alexsetta/broker/cfg"
	"github.com/alexsetta/broker/cotacao"
	"github.com/alexsetta/broker/tipos"
	"github.com/alexsetta/broker/util"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	hora     = time.Now().Add(time.Hour * -5)
	alerta   = tipos.Alertas{hora, hora, hora, hora, hora, hora}
	carteira = tipos.Carteira{}
	config   = tipos.Config{}
	loc      = time.FixedZone("UTC-3", -3*60*60)
)

func PrettyJson(data interface{}) (string, error) {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")

	err := encoder.Encode(data)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func test1() {
	var filename string = "./coletor.log"
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	Formatter := new(log.JSONFormatter)
	Formatter.TimestampFormat = "2006-01-02 15:04:05"
	log.SetFormatter(Formatter)
	if err != nil {
		fmt.Println(err)
	} else {
		log.SetOutput(f)
	}
	log.Info("inicio")

	if err := cfg.ReadConfig("d:/dev/go/app/coletor/carteira.cfg", &carteira); err != nil {
		log.Fatal(err)
	}
	resp, _, result, err := cotacao.Calculo(carteira.Ativos[0], config, alerta, nil)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(resp)

	prettyJSON, err := PrettyJson(result)
	if err != nil {
		log.Fatal("Failed to generate json", err)
	}
	fmt.Printf("%s\n", prettyJSON)

	rsi, err := cotacao.GetRSI("https://br.investing.com/crypto/bitcoin/btc-brl-technical")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	fmt.Println(rsi)
}

func test2() {
	res, err := http.Get("https://br.investing.com/crypto/cardano/ada-brl")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	doc := string(b)
	fmt.Println(doc)
}

func test3() {
	s, err := util.GetHttp("https://api.binance.com/api/v3/ticker/price?symbol=ADABRL")
	util.FailOnError(err, "Failed to get http")
	fmt.Println(s)
}

func test4() {
	s := "12.17400000"
	v, err := strconv.ParseFloat(s, 64)
	util.FailOnError(err, "Failed to convert string to float")
	fmt.Println(v)
}

func test5() {
	fmt.Println(util.USDToBRL(10))
}

func slicetostring(slice []string) string {
	var str string
	for _, v := range slice {
		str += v + ","
	}
	return str[:len(str)-1]
}

func test6() {
	execs := []string{"00", "20", "40"}
	fmt.Println("Horário para execução: " + slicetostring(execs))
	for _, exec := range execs {
		fmt.Println(exec)
	}
}

func main() {
	test6()
}
