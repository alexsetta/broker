package loga

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"golang.org/x/exp/slog"
	"os"
)

func NewLog(filename string) {
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
	} else {
		log.SetOutput(f)
		slog.SetDefault(slog.New(slog.NewTextHandler(f)))
	}
}
