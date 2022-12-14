package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func Ticker(c *gin.Context) {
	id := c.Param("id")

	asset, err := NewAsset(id, true)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	log.Println("Ticker", id)

	resposta := ""
	if id == "all" {
		data, err := asset.GetAll()
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		resposta, _ = PrettyJson(data)
		c.String(http.StatusOK, resposta)
		return
	}

	err = asset.Find()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	prettyJSON, _ := PrettyJson(asset.data)
	resposta = prettyJSON + ","
	resposta = "[" + resposta[:len(resposta)-1] + "]"
	c.String(http.StatusOK, resposta)
}
