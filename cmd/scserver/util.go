package main

import (
	"bytes"
	"encoding/json"
	"github.com/alexsetta/broker/pkg/cfg"
	"net/http"
)

func formatJSONError(message string) []byte {
	appError := struct {
		Message string `json:"message"`
	}{
		message,
	}
	response, err := json.Marshal(appError)
	if err != nil {
		return []byte(err.Error())
	}
	return response
}

func PrettyJson(data interface{}) (string, error) {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")

	err := encoder.Encode(data)
	if err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func saidaJson(w http.ResponseWriter, outJson interface{}) {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(outJson)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(formatJSONError("Erro convertendo em JSON"))
		return
	}
}

func saidaHtml(w http.ResponseWriter, p Page) {
	if err := templates.ExecuteTemplate(w, "result.html", p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func ReadConfig() error {
	dirBase := "../../config/"
	if err := cfg.ReadConfig(dirBase+"broker.cfg", &config); err != nil {
		return err
	}

	if err := cfg.ReadConfig(dirBase+"wallet.cfg", &carteira); err != nil {
		return err
	}

	return nil
}
