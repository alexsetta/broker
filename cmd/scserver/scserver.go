package main

import (
	"crypto/tls"
	"github.com/alexsetta/broker/pkg/tipos"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
)

// Generated by https://quicktype.io
type Page struct {
	Result string
}

var (
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	templates *template.Template
	porta     string
	hora      = time.Now().Add(time.Hour * -5)
	alerta    = tipos.Alertas{hora, hora, hora, hora, hora, hora}
	carteira  = tipos.Carteira{}
	config    = tipos.Config{}
	start     = time.Now()
)

const (
	dirBase   = "../.."
	dirFiles  = dirBase + "/files/"
	dirConfig = dirBase + "/config/"
)

func main() {
	porta = "8081"
	if len(os.Args) == 2 {
		porta = os.Args[1]
	}
	templates = template.Must(template.ParseFiles("./templates/result.html"))

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/", Index)
	r.GET("/show/:from/:to", Order)
	r.GET("/order/:from/:to", Order)
	r.GET("/total", Total)
	r.GET("/ticker/:id", Ticker)

	log.Println("Listen port " + porta)
	log.Fatal(r.Run(":" + porta))
}

func Index(c *gin.Context) {
	c.String(http.StatusOK, "scserver online\nstart time: %s", time.Since(start))
}
