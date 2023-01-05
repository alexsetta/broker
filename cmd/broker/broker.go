package main

import (
	"fmt"
	"github.com/alexsetta/broker/pkg/cfg"
	"github.com/alexsetta/broker/pkg/common"
	"github.com/alexsetta/broker/pkg/loga"
	"github.com/alexsetta/broker/pkg/mensagem"
	"github.com/alexsetta/broker/pkg/price"
	"github.com/alexsetta/broker/pkg/rsi"
	"github.com/alexsetta/broker/pkg/tipos"
	"github.com/alexsetta/broker/pkg/util"
	log "github.com/sirupsen/logrus"
	"os"
	"time"
)

var (
	hora      = time.Now().Add(time.Hour * -5)
	alerta    = tipos.Alertas{hora, hora, hora, hora, hora, hora}
	carteira  = tipos.Carteira{}
	config    = tipos.Config{}
	loc       = time.FixedZone("UTC-3", -3*60*60)
	ultimoDia = "00"
	dir       = common.NewDir("")
	file      = common.NewFile(dir)
)

func main() {
	loga.NewLog(file.Log)
	log.Info("Broker iniciado")

	if err := cfg.ReadConfig(file.Config, &config); err != nil {
		log.Fatal(err)
	}
	if config.TimeNextAlert == 0.0 {
		config.TimeNextAlert = 1
	}

	// cria um RSI para cada ativo
	if err := cfg.ReadConfig(file.Wallet, &carteira); err != nil {
		log.Fatal(err)
	}
	mr := make(map[string]*rsi.RSI)
	for _, atv := range carteira.Ativos {
		mr[atv.Simbolo] = rsi.NewRSI(atv.Simbolo)
		//	if len(atv.RSI) > 0 {
		//		mr[atv.Simbolo].Load()
		//	}
	}

	for {
		if err := cfg.ReadConfig(file.Wallet, &carteira); err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}
		hm := time.Now().In(loc).Format("15:04")
		header := fmt.Sprintln(util.Now())
		ultimo := header
		fmt.Print("\n" + header)

		for _, ativo := range carteira.Ativos {
			go func(ativo tipos.Ativo, cfg tipos.Config, alerta tipos.Alertas) {
				resp, semaforo, _, err := price.Get(ativo, config, alerta, mr)
				if err != nil {
					log.Println(err)
					return
				}
				fmt.Println(resp)
				ultimo += resp
				setAlert(semaforo)
			}(ativo, config, alerta)
		}
		if err := os.WriteFile(file.Last, []byte(ultimo), 0644); err != nil {
			fmt.Println(fmt.Errorf("writefile: %w", err))
		}

		dia := time.Now().In(loc).Format("02")
		if hm >= "07:00" && hm <= "08:00" && dia != ultimoDia {
			ultimoDia = dia
			if err := mensagem.Send(config, ultimo); err != nil {
				fmt.Println(fmt.Errorf("send: %w", err))
			}
		}
		time.Sleep(time.Duration(config.SleepSeconds) * time.Second)
	}
}

func setAlert(tipo string) bool {
	switch tipo {
	case "ganho":
		if time.Since(alerta.Ganho).Hours() < config.TimeNextAlert {
			return false
		}
		alerta.Ganho = time.Now()
	case "perda":
		if time.Since(alerta.Perda).Hours() < config.TimeNextAlert {
			return false
		}
		alerta.Perda = time.Now()
	case "alertainf":
		if time.Since(alerta.AlertaInf).Hours() < config.TimeNextAlert {
			return false
		}
		alerta.AlertaInf = time.Now()
	case "alertasup":
		if time.Since(alerta.AlertaSup).Hours() < config.TimeNextAlert {
			return false
		}
		alerta.AlertaSup = time.Now()
	case "alertaperc":
		if time.Since(alerta.AlertaPerc).Hours() < config.TimeNextAlert {
			return false
		}
		alerta.AlertaPerc = time.Now()
	case "rsi":
		if time.Since(alerta.RSI).Hours() < config.TimeNextAlert {
			return false
		}
		alerta.RSI = time.Now()
	}
	return true
}
