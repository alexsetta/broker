package main

import (
	"fmt"
	"github.com/alexsetta/broker/cfg"
	"github.com/alexsetta/broker/cotacao"
	"github.com/alexsetta/broker/rsi"
	"net/http"
	"time"
)

func Total(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Total")
	if err := cfg.ReadConfig("../coletor.cfg", &config); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// desabilita mensagens no Telegram
	config.TelegramID = 0

	if err := cfg.ReadConfig("../carteira.cfg", &carteira); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mr := make(map[string]*rsi.RSI)
	atual := 0.0
	for _, atv := range carteira.Ativos {
		if atv.Tipo != "criptomoeda" {
			continue
		}
		mr[atv.Simbolo] = rsi.NewRSI(atv.Simbolo, "../../files", false)

		_, _, out, err := cotacao.Calculo(atv, config, alerta, mr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		atual += out.Atual
	}

	resposta := fmt.Sprintf(`{"hora": "%v","total": %v}`, time.Now().In(time.FixedZone("UTC-3", -3*60*60)).Format("02/01/2006 15:04:05"), atual)
	w.Write([]byte(resposta))

}
