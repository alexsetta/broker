package tipos

import (
	"time"
)

type Config struct {
	SleepSeconds  int     `json:"sleepSeconds"`
	TimeNextAlert float64 `json:"timeNextAlert"`
	TelegramID    int64   `json:"telegramID"`
	TelegramToken string  `json:"telegramToken"`
	EmailLogin    string  `json:"emailLogin"`
	EmailPassword string  `json:"emailPassword"`
	EmailTo       string  `json:"emailTo"`
	SaveLog       bool    `json:"saveLog"`
}

type Quotes struct {
	Quotes []Quote `json:"quotes"`
}

type Quote struct {
	Symbol   string  `json:"symbol"`
	Quantity int64   `json:"quantity"`
	Stop     float64 `json:"stop"`
}

type Carteira struct {
	Ativos []Ativo `json:"ativos"`
}

type Ativo struct {
	Simbolo    string  `json:"simbolo"`
	Link       string  `json:"link"`
	RSI        string  `json:"rsi"`
	Tipo       string  `json:"tipo"`
	Taxa       float64 `json:"taxa"`
	Quantidade float64 `json:"quantidade"`
	Inicial    float64 `json:"inicial"`
	Perda      float64 `json:"perda"`
	Ganho      float64 `json:"ganho"`
	AlertaInf  float64 `json:"alerta_inf"`
	AlertaSup  float64 `json:"alerta_sup"`
	AlertaPerc float64 `json:"alerta_perc"`
}

type Alertas struct {
	Ganho      time.Time
	Perda      time.Time
	AlertaInf  time.Time
	AlertaSup  time.Time
	AlertaPerc time.Time
	RSI        time.Time
}

type Result struct {
	Hora               string  `json:"hora"`
	Simbolo            string  `json:"simbolo"`
	Quantidade         float64 `json:"quantidade"`
	Inicial            float64 `json:"inicial"`
	Preco              float64 `json:"preco"`
	Resultado          float64 `json:"resultado"`
	Moeda              string  `json:"moeda"`
	Atual              float64 `json:"atual"`
	Percentual         float64 `json:"percentual"`
	RSI                float64 `json:"rsi"`
	PriceChange        float64 `json:"price_change"`
	PriceChangePercent float64 `json:"priceChangePercent"`
	LastQty            float64 `json:"last_qty"`
	Volume             float64 `json:"volume"`
}

var Simula bool
