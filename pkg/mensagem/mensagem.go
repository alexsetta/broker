package mensagem

import (
	"crypto/tls"
	"fmt"
	"github.com/alexsetta/broker/pkg/tipos"
	"github.com/alexsetta/telegram"
	"gopkg.in/gomail.v2"
)

func Send(cfg tipos.Config, msg string) error {
	var err error
	errTelegram := sendTelegram(cfg, msg)
	errEmail := sendEmail(cfg, msg)

	if errTelegram != nil {
		err = errTelegram
	}

	if errEmail != nil {
		err = errEmail
	}

	return err
}

func sendTelegram(cfg tipos.Config, msg string) error {
	if cfg.TelegramID == 0 || cfg.TelegramToken == "" {
		return nil
	}

	if err := telegram.SendMessage(cfg.TelegramID, cfg.TelegramToken, msg); err != nil {
		return fmt.Errorf("sendTelegram: %w", err)
	}
	return nil
}

func sendEmail(cfg tipos.Config, msg string) error {
	if cfg.EmailTo == "" || cfg.EmailLogin == "" || cfg.EmailPassword == "" {
		return nil
	}
	m := gomail.NewMessage()
	m.SetHeader("From", cfg.EmailLogin)
	m.SetHeader("To", cfg.EmailTo)
	m.SetHeader("Subject", "SMARTBOT")
	m.SetBody("text/html", msg)

	d := gomail.NewDialer("smtp.gmail.com", 587, cfg.EmailLogin, cfg.EmailPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("sendEmail: %w", err)
	}
	return nil
}
