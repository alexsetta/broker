package main

import (
	"fmt"
	"github.com/alexsetta/broker/pkg/cfg"
	"github.com/alexsetta/broker/pkg/tipos"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	dirBase = "../.."
)

var (
	carteira = tipos.Carteira{}
	config   = tipos.Config{}
)

func main() {
	if err := cfg.ReadConfig(dirBase+"/config/broker.cfg", &config); err != nil {
		log.Fatal(fmt.Sprintf("price: read broker.cfg: %s", err))
	}
	config.TelegramID = 0

	bot, err := tgbotapi.NewBotAPI(config.TelegramToken)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, _ := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			var response string
			cmd := strings.ToLower(update.Message.Text)
			switch {
			case cmd == "status":
				response = "online"
			case cmd == "total":
				response = Total()
			case cmd == "eth":
				response = Cotacao("ethbrl")
			case cmd == "btc":
				response = Cotacao("btcbrl")
			case cmd[0:1] == "/":
				response = Cotacao(cmd[1:])
			default:
				response = "Comando inválido: " + cmd
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
			bot.Send(msg)
		}
	}
}
