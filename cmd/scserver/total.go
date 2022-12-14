package main

import (
	"fmt"
	"github.com/alexsetta/broker/pkg/cfg"
	"github.com/alexsetta/broker/pkg/price"
	"github.com/alexsetta/broker/pkg/rsi"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Total(c *gin.Context) {
	if err := cfg.ReadConfig(dirConfig+"broker.cfg", &config); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	// desabilita mensagens no Telegram
	config.TelegramID = 0

	if err := cfg.ReadConfig(dirConfig+"wallet.cfg", &carteira); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	mr := make(map[string]*rsi.RSI)
	atual := 0.0
	for _, atv := range carteira.Ativos {
		if atv.Tipo != "criptomoeda" {
			continue
		}
		mr[atv.Simbolo] = rsi.NewRSI(atv.Simbolo)

		_, _, out, err := price.Get(atv, config, alerta, mr)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		atual += out.Atual
	}

	resposta := fmt.Sprintf(`{"hora": "%v","total": %v}`, time.Now().In(time.FixedZone("UTC-3", -3*60*60)).Format("02/01/2006 15:04:05"), atual)
	c.String(http.StatusOK, resposta)
}
