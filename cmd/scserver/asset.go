package main

import (
	"errors"
	"fmt"
	"github.com/alexsetta/broker/cotacao"
	"github.com/alexsetta/broker/rsi"
	"github.com/alexsetta/broker/tipos"
	"strings"
)

type Asset struct {
	id       string
	loadFile bool
	data     tipos.Result
}

func NewAsset(id string, loadFile bool) (*Asset, error) {
	asset := &Asset{id: id, loadFile: loadFile, data: tipos.Result{}}

	if err := asset.IsValid(); err != nil {
		return nil, err
	}

	if id == "all" {
		return asset, nil
	}

	if err := asset.Find(); err != nil {
		return nil, err
	}
	return asset, nil
}

func (a *Asset) SetLoadFile(loadFile bool) {
	a.loadFile = loadFile
}

func (a *Asset) IsValid() error {
	if a.id == "" {
		return errors.New("id is empty")
	}
	return nil
}

func (a *Asset) Find() error {
	if err := ReadConfig(); err != nil {
		return err
	}
	// desabilita mensagens no Telegram
	config.TelegramID = 0

	// desabilita saveLog
	config.SaveLog = false
	mr := make(map[string]*rsi.RSI)

	ativo := tipos.Ativo{}
	for _, atv := range carteira.Ativos {
		mr[atv.Simbolo] = rsi.NewRSI(atv.Simbolo, dirPath, a.loadFile)
		if strings.ToLower(atv.Simbolo) == a.id {
			if a.loadFile && len(atv.RSI) > 0 {
				mr[atv.Simbolo].Load()
			}
			ativo = atv
			break
		}
	}

	if ativo == (tipos.Ativo{}) {
		return fmt.Errorf("asset %s not found", a.id)
	}

	_, _, out, err := cotacao.Calculo(ativo, config, alerta, mr)
	if err != nil {
		return err
	}
	a.data = out
	return nil
}

func (a *Asset) GetAll() ([]tipos.Result, error) {
	if err := ReadConfig(); err != nil {
		return []tipos.Result{}, err
	}
	// desabilita mensagens no Telegram
	config.TelegramID = 0

	// desabilita saveLog
	config.SaveLog = false

	mr := make(map[string]*rsi.RSI)

	resposta := ""
	var outJson []tipos.Result
	for _, atv := range carteira.Ativos {
		mr[atv.Simbolo] = rsi.NewRSI(atv.Simbolo, dirPath, a.loadFile)
		_, _, out, err := cotacao.Calculo(atv, config, alerta, mr)
		if err != nil {
			return []tipos.Result{}, err
		}
		outJson = append(outJson, out)
		prettyJSON, _ := PrettyJson(out)
		resposta += prettyJSON + ","
	}

	return outJson, nil
}
