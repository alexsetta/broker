package main

import (
	"flag"
	"fmt"
	"github.com/alexsetta/broker/pkg/cfg"
	"github.com/alexsetta/broker/pkg/common"
	"github.com/alexsetta/broker/pkg/cotacao"
	"github.com/alexsetta/broker/pkg/loga"
	"github.com/alexsetta/broker/pkg/mensagem"
	"github.com/alexsetta/broker/pkg/rsi"
	"github.com/alexsetta/broker/pkg/tipos"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	hora      = time.Now().Add(time.Hour * -5)
	alerta    = tipos.Alertas{hora, hora, hora, hora, hora, hora}
	carteira  = tipos.Carteira{}
	config    = tipos.Config{}
	loc       = time.FixedZone("UTC-3", -3*60*60)
	ultimoDia = "00"
	dir       = common.NewDir("")
	filename  = dir.Files + "ultimo.txt"
)

func main() {
	loga.NewLog(dir.Files + "monitor.log")
	log.Println("Monitor iniciado")

	test := flag.Bool("test", false, "teste Telegram")
	simula := flag.Bool("simula", false, "simula trade")
	flag.Parse()

	tipos.Simula = *simula
	if tipos.Simula {
		fmt.Println("Simulando trade...")
	}

	if err := cfg.ReadConfig(dir.Config+"/broker.cfg", &config); err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	if config.TimeNextAlert == 0.0 {
		config.TimeNextAlert = 1
	}

	if *test {
		if err := mensagem.Send(config, "Mensagem de teste"); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	// cria um RSI para cada ativo
	if err := cfg.ReadConfig(dir.Config+"carteira.cfg", &carteira); err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	mr := make(map[string]*rsi.RSI)
	for _, atv := range carteira.Ativos {
		mr[atv.Simbolo] = rsi.NewRSI(atv.Simbolo, dir.Files, false)
		if len(atv.RSI) > 0 {
			mr[atv.Simbolo].Load()
		}
	}

	for {
		if err := cfg.ReadConfig(dir.Config+"carteira.cfg", &carteira); err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}
		hm := time.Now().In(loc).Format("15:04")
		header := fmt.Sprintln(time.Now().In(loc).Format("02/01/2006 15:04:05"))
		ultimo := header
		fmt.Print("\n" + header)

		for _, ativo := range carteira.Ativos {
			go func(ativo tipos.Ativo, cfg tipos.Config, alerta tipos.Alertas) {
				resp, semaforo, _, err := cotacao.Calculo(ativo, config, alerta, mr)
				if err != nil {
					log.Println(err)
					return
				}
				fmt.Println(resp)
				ultimo += resp
				setAlert(semaforo)
			}(ativo, config, alerta)
		}
		if err := os.WriteFile(filename, []byte(ultimo), 0644); err != nil {
			fmt.Println(fmt.Errorf("writefile: %w", err))
		}

		//fmt.Println()
		dia := time.Now().In(loc).Format("02")
		if hm >= "07:00" && hm <= "08:00" && dia != ultimoDia {
			ultimoDia = dia
			fmt.Println("Mensagem diária")
			msg_diaria := total(config, carteira)
			if err := mensagem.Send(config, msg_diaria); err != nil {
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

func total(cfg tipos.Config, cart tipos.Carteira) string {
	atual := 0.0
	mr := make(map[string]*rsi.RSI)
	for _, atv := range cart.Ativos {
		if atv.Tipo != "criptomoeda" {
			continue
		}
		mr[atv.Simbolo] = rsi.NewRSI(atv.Simbolo, dir.Files, false)
		_, _, out, err := cotacao.Calculo(atv, cfg, alerta, mr)
		if err != nil {
			return fmt.Sprintf("total: %w", err)
		}
		atual += out.Atual
	}

	return fmt.Sprintf(`{"hora": "%v","total": %v}`, time.Now().In(time.FixedZone("UTC-3", -3*60*60)).Format("02/01/2006 15:04:05"), atual)
}
