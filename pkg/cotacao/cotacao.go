package cotacao

import (
	"fmt"
	"github.com/alexsetta/broker/pkg/mensagem"
	"github.com/alexsetta/broker/pkg/rsi"
	"github.com/alexsetta/broker/pkg/tipos"
	"github.com/alexsetta/broker/pkg/util"
	"io"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	re    = regexp.MustCompile(`<span class="instrument-price_last__KQzyA" data-test="instrument-price-last">(.*?)</span>`)
	reTV  = regexp.MustCompile(`<span class="arial_26 inlineblock pid-(\d*)-last" id="last_last" dir="ltr">(.*?)</span>`)
	reBNC = regexp.MustCompile(`"([^"]+)"`)

	clientHttp = &http.Client{
		Timeout: time.Second * 15,
	}
)

func getHttp(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("gethttp: %w", err)
	}
	req.Header.Add("User-Agent", "XYZ/3.0")
	resp, err := clientHttp.Do(req)
	if err != nil {
		return "", fmt.Errorf("gethttp: %w", err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("gethttp: %w", err)
	}
	return string(b), nil
}

func Calculo(ativo tipos.Ativo, cfg tipos.Config, alerta tipos.Alertas, rsi map[string]*rsi.RSI) (string, string, tipos.Result, error) {
	var result tipos.Result
	price, m, err := Price(ativo)
	if err != nil {
		return "", "", result, fmt.Errorf("calculo: %w", err)
	}

	moeda := "R$"
	dolar := 1.0
	if strings.Contains(ativo.Simbolo, "USD") && ativo.Simbolo != "USD" {
		moeda = "US$"
		//dolar = util.USDToBRL(dolar)
	}
	price *= dolar

	sema := ""
	taxa := 1.00 - ativo.Taxa
	base := ativo.Quantidade * ativo.Inicial * dolar
	atual := ativo.Quantidade * price * taxa
	diff := atual - base
	perc := ((price * 100 / (ativo.Inicial * dolar)) - 100) * taxa

	result.Hora = time.Now().In(time.FixedZone("UTC-3", -3*60*60)).Format("02/01/2006 15:04:05")
	result.Simbolo = ativo.Simbolo
	result.Quantidade = ativo.Quantidade
	result.Inicial = ativo.Inicial
	result.Atual = math.Trunc(atual*100) / 100
	result.Resultado = math.Trunc(diff*100) / 100
	result.Preco = math.Trunc(price*100) / 100
	result.RSI = 0.00
	result.Percentual = math.Trunc(perc*100) / 100

	if (ativo.Tipo == "criptomoeda" || ativo.Tipo == "etf") && len(ativo.RSI) > 0 {
		//result.RSI, _ = GetRSI(ativo.RSI)
		rsi[ativo.Simbolo].Add(result.Preco)
		result.RSI = rsi[ativo.Simbolo].Calculate()
	}

	if len(m) >= 29 {
		result.PriceChange = util.StringToValue(m[3])
		result.PriceChangePercent = util.StringToValue(m[5])
		result.LastQty = util.StringToValue(m[13])
		result.Volume = util.StringToValue(m[29])
	}
	if cfg.SaveLog {
		util.AppendFile("./cotacao.log", fmt.Sprintf("%v;%.8f;%.2f;%.9f;%0.0f;%0.0f", ativo.Simbolo, price, result.RSI, result.PriceChange, result.Volume, result.LastQty))
	}

	res := fmt.Sprintf("%-12v %-20v %-15v", ativo.Simbolo, fmt.Sprintf("(%s) Preço: %.2f", moeda, price), fmt.Sprintf("Dif.: %.2f", diff))
	if ativo.AlertaPerc != 0 || ativo.Tipo == "criptomoeda" {
		rsiCalc := result.RSI
		res += fmt.Sprintf("%-10v %-22v", fmt.Sprintf(" (%.2f%%)", perc), fmt.Sprintf("Total: %.2f ", atual))
		if len(ativo.RSI) > 0 {
			res += fmt.Sprintf("%-12v", fmt.Sprintf("RSI: %.2f", rsiCalc))
		}
		if rsiCalc != 0 && (rsiCalc <= 30 || rsiCalc >= 70) && time.Since(alerta.RSI).Hours() > cfg.TimeNextAlert {
			acao := "VENDA"
			if rsiCalc <= 30 {
				acao = "COMPRA"
			}
			msg := res + "RSI Alerta de " + acao

			sema := "rsi"
			_ = mensagem.Send(cfg, msg)
			return msg, sema, result, nil
		}
	}

	if ativo.Perda != 0 && diff < 0 && diff < ativo.Perda && time.Since(alerta.Perda).Hours() > cfg.TimeNextAlert {
		msg := res + "Atingiu o limite de perda!"
		sema := "perda"
		_ = mensagem.Send(cfg, msg)
		return msg, sema, result, nil
	}
	if ativo.Ganho != 0 && diff > 0 && diff > ativo.Ganho && time.Since(alerta.Ganho).Hours() > cfg.TimeNextAlert {
		msg := res + "Atingiu o limite de ganho!"
		sema := "ganho"
		_ = mensagem.Send(cfg, msg)
		return msg, sema, result, nil
	}

	if ativo.AlertaInf != 0 && price <= ativo.AlertaInf && time.Since(alerta.AlertaInf).Hours() > cfg.TimeNextAlert {
		msg := res + "Atingiu o limite inferior!"
		sema := "alertainf"
		_ = mensagem.Send(cfg, msg)
		return msg, sema, result, nil
	}

	if ativo.AlertaSup != 0 && price >= ativo.AlertaSup && time.Since(alerta.AlertaSup).Hours() > cfg.TimeNextAlert {
		msg := res + "Atingiu o limite superior!"
		sema := "alertasup"
		_ = mensagem.Send(cfg, msg)
		return msg, sema, result, nil
	}

	if ativo.AlertaPerc > 0 && perc > ativo.AlertaPerc && time.Since(alerta.AlertaPerc).Hours() > cfg.TimeNextAlert {
		msg := res + "Atingiu o limite percentual!"
		sema := "alertaperc"
		_ = mensagem.Send(cfg, msg)
		return msg, sema, result, nil
	}

	return res, sema, result, nil
}

func GetRSI(url string) (float64, error) {
	doc, err := getHttp(url)
	if err != nil {
		return 0, fmt.Errorf("valor: %w", err)
	}
	reRSI := regexp.MustCompile(`<td class="right">(.*?)</td>`)
	matches := reRSI.FindStringSubmatch(doc)
	if len(matches) != 2 {
		return 0, fmt.Errorf("getRSI: valor do RSI não encontrado")
	}

	s := matches[1]
	ponto := strings.Index(s, ".")
	virgula := strings.Index(s, ",")

	if ponto < virgula {
		if strings.Contains(s, ",") {
			s = strings.ReplaceAll(s, ".", "")
		}
	}
	s = strings.ReplaceAll(s, ",", ".")
	price, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("valor: %w", err)
	}

	return price, nil
}

func priceString(ativo tipos.Ativo, doc string) (string, []string, error) {
	var matches []string
	var s string

	if strings.Contains(ativo.Link, "binance") {
		matches = reBNC.FindAllString(doc, 32)
		if len(matches) < 11 {
			return "", matches, fmt.Errorf("priceString: cotação não encontrada: %s", ativo.Simbolo)
		}
		s = strings.ReplaceAll(matches[11], `"`, "")
		s = strings.ReplaceAll(s, ".", ",")
		return s, matches, nil
	}

	if ativo.Tipo == "acao" {
		matches = re.FindStringSubmatch(doc)
		if len(matches) != 2 {
			return "", matches, fmt.Errorf("priceString: cotação não encontrada: %s", ativo.Simbolo)
		}
		s = matches[1]
	} else if ativo.Tipo == "etf" {
		matches = reTV.FindStringSubmatch(doc)
		if len(matches) != 3 {
			return "", matches, fmt.Errorf("priceString: cotação não encontrada: %s", ativo.Simbolo)
		}
		s = matches[2]
	} else {
		matches = reTV.FindStringSubmatch(doc)
		if len(matches) != 3 {
			return "", matches, fmt.Errorf("priceString: cotação não encontrada: %s", ativo.Simbolo)
		}
		s = matches[2]
	}
	return s, matches, nil
}

func Price(ativo tipos.Ativo) (float64, []string, error) {
	var s string
	var m []string

	if ativo.Simbolo == "USD" {
		brl, err := util.USD()
		if err != nil {
			return 0, m, fmt.Errorf("price: %w", err)
		}
		return brl, m, nil
	}

	doc, err := getHttp(ativo.Link)
	if err != nil {
		return 0, m, fmt.Errorf("price: %w", err)
	}

	s, m, err = priceString(ativo, doc)
	if err != nil {
		return 0, m, fmt.Errorf("price: %w", err)
	}

	ponto := strings.Index(s, ".")
	virgula := strings.Index(s, ",")
	//fmt.Printf("DEBUG: %s %d %d\n", s, ponto, virgula)
	if virgula < 0 {
		s = s + ",00"
		virgula = strings.Index(s, ",")
	}
	if ponto < virgula {
		if strings.Contains(s, ",") {
			s = strings.ReplaceAll(s, ".", "")
		}
	}
	s = strings.ReplaceAll(s, ",", ".")
	price, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, m, fmt.Errorf("valor: %w", err)
	}

	return price, m, nil
}
