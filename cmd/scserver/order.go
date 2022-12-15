package main

import (
	"fmt"
	"github.com/alexsetta/broker/util"
	"github.com/gin-gonic/gin"
	"os"
	"strings"

	"log"
	"net/http"
)

//type AssetOrder struct {
//	qty   float64
//	value float64
//}

type OrderResult struct {
	Time      string
	From      string
	FromQty   float64
	FromValue float64
	To        string
	ToQty     float64
	ToValue   float64
}

func (r *OrderResult) String() string {
	res := fmt.Sprintf(`{"time":"%s","from":"%s","fromQty":%f,"fromValue":%f,"to":"%s","toQty":%f,"toValue":%f}`, r.Time, r.From, r.FromQty, r.FromValue, r.To, r.ToQty, r.ToValue)
	util.AppendFile(dirFiles+"order.log", res)
	return res
}

func (r *OrderResult) Json() string {
	return fmt.Sprintf(`{"time":"%s","from":"%s","fromQty":%f,"fromValue":%f,"to":"%s","toQty":%f,"toValue":%f}`, r.Time, r.From, r.FromQty, r.FromValue, r.To, r.ToQty, r.ToValue)
}

func (r *OrderResult) Save() error {
	if err := os.WriteFile(dirFiles+"order.json", []byte(r.Json()), 0644); err != nil {
		return err
	}
	return nil
}

func Order(c *gin.Context) {
	from := c.Param("from")
	if from == "" {
		c.String(http.StatusBadRequest, "argument 'from' is required")
		return
	}

	to := c.Param("to")
	if to == "" {
		c.String(http.StatusBadRequest, "argument 'to' is required")
		return
	}

	verb := strings.Split(c.Request.URL.Path, "/")[1]
	log.Println(verb+" from:", from, "to:", to)

	fromAsset, err := NewAsset(from, false)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	toAsset, err := NewAsset(to, false)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	res := OrderResult{
		Time:      util.Now(),
		From:      fromAsset.data.Simbolo,
		FromValue: fromAsset.data.Preco,
		FromQty:   fromAsset.data.Quantidade,
		To:        toAsset.data.Simbolo,
		ToValue:   toAsset.data.Preco,
		ToQty:     fromAsset.data.Quantidade / toAsset.data.Preco,
	}

	if verb == "order" {
		_ = res.Save()
	}
	fmt.Println(res.String())
	s, err := PrettyJson(res)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, s)
}
