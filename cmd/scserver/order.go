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
	time      string
	from      string
	fromQty   float64
	fromValue float64
	to        string
	toQty     float64
	toValue   float64
}

func (r *OrderResult) String() string {
	res := fmt.Sprintf(`{"time":"%s","from":"%s","fromQty":%f,"fromValue":%f,"to":"%s","toQty":%f,"toValue":%f}`, r.time, r.from, r.fromQty, r.fromValue, r.to, r.toQty, r.toValue)
	util.AppendFile(dirFiles+"order.log", res)
	return res
}

func (r *OrderResult) Json() string {
	return fmt.Sprintf(`{"time":"%s","from":"%s","fromQty":%f,"fromValue":%f,"to":"%s","toQty":%f,"toValue":%f}`, r.time, r.from, r.fromQty, r.fromValue, r.to, r.toQty, r.toValue)
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
		time:      util.Now(),
		from:      fromAsset.data.Simbolo,
		fromValue: fromAsset.data.Preco,
		fromQty:   fromAsset.data.Quantidade,
		to:        toAsset.data.Simbolo,
		toValue:   toAsset.data.Preco,
		toQty:     fromAsset.data.Quantidade / toAsset.data.Preco,
	}

	if verb == "order" {
		_ = res.Save()
	}
	fmt.Println(res.String())
	c.String(http.StatusOK, res.String())
}
