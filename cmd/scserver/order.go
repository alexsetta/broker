package main

import (
	"fmt"
	"github.com/alexsetta/broker/util"
	"os"

	//"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
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
	return fmt.Sprintf(`{"time":"%s","from":"%s","fromQty":%f,"fromValue":%f,"to":"%s","toQty":%f,"toValue":%f}`,
		r.time, r.from, r.fromQty, r.fromValue, r.to, r.toQty, r.toValue)
}

func (r *OrderResult) Save() error {
	if err := os.WriteFile(dirFiles+"order.json", []byte(r.Json()), 0644); err != nil {
		return err
	}
	return nil
}

func Order(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	from := vars["from"]
	if from == "" {
		http.Error(w, "Argument 'from' not found", http.StatusInternalServerError)
		return
	}

	to := vars["to"]
	if to == "" {
		http.Error(w, "Argument 'to' not found", http.StatusInternalServerError)
		return
	}
	log.Println("Order from", from, "to", to)

	fromAsset, err := NewAsset(from, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	toAsset, err := NewAsset(to, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	_ = res.Save()
	fmt.Println(res.String())
	fmt.Fprint(w, res.Json())
}
