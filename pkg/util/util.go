package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func Now() string {
	return time.Now().In(time.FixedZone("UTC-3", -3*60*60)).Format("2006/01/02 15:04:05")
}

func Segundos(s string) (int64, error) {
	v := strings.Split(s, ":")

	var i int64 = 0
	var t int64 = 0
	var err error

	i, err = strconv.ParseInt(v[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("segundos: %w", err)
	}
	t += i * 3600

	i, err = strconv.ParseInt(v[1], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("segundos: %w", err)
	}
	t += i * 60

	i, err = strconv.ParseInt(v[2], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("segundos: %w", err)
	}
	t += i

	return t, nil
}

func Find(slice []string, val string) int {
	for i, item := range slice {
		if item == val {
			return i
		}
	}
	return -1
}

func Contains(slice []string, val string) bool {
	for _, item := range slice {
		if val == item {
			return true
		}
	}
	return false
}

func RotateLog(file string, size int64, backups int) {
	fi, err := os.Stat(file)
	if err != nil {
		return
	}
	if fi.Size() < size {
		return
	}

	f := fmt.Sprintf("%v.%v", file, backups)
	_ = os.Remove(f)
	for i := backups; i > 1; i-- {
		fs := fmt.Sprintf("%v.%v", file, i)
		fi := fmt.Sprintf("%v.%v", file, i-1)
		_ = os.Rename(fi, fs)
	}
	f = fmt.Sprintf("%v.%v", file, 1)
	_ = os.Rename(file, f)
}

func Debug(msg string) {
	hourExec := time.Now().In(time.FixedZone("UTC-3", -3*60*60)).Format("2006/01/02 15:04:05")
	fmt.Println(hourExec + " " + msg)
}

func StringToValue(s string) float64 {
	s = strings.ReplaceAll(s, `"`, "")
	s = strings.ReplaceAll(s, `'`, "")
	s = strings.ReplaceAll(s, `,`, "")
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return -1.00
	}
	return f
}

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func AppendFile(filename, message string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(fmt.Errorf("appendfile: openfile: %w", err))
	}
	defer file.Close()
	message = fmt.Sprint(time.Now().In(time.FixedZone("UTC-3", -3*60*60)).Format("02/01/2006 15:04:05")) + ";" + message + "\n"
	if _, err := file.WriteString(message); err != nil {
		fmt.Println(fmt.Errorf("appendfile: writestring: %w", err))
	}
}

func GetHttp(url string) (string, error) {
	clientHttp := &http.Client{
		Timeout: time.Second * 15,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("gethttp: %w", err)
	}
	req.Header.Add("User-Agent", "XYZ/3.0")
	resp, err := clientHttp.Do(req)
	if err != nil {
		return "", fmt.Errorf("gethttp: %w", err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("gethttp: %w", err)
	}
	return string(b), nil
}

func USD() (float64, error) {
	s, err := GetHttp("https://economia.awesomeapi.com.br/json/last/usd-brl")
	if err != nil {
		return 0, fmt.Errorf("usd: gethttp: %w", err)
	}

	type result struct {
		USDBRL struct {
			Price string `json:"bid"`
		} `json:"USDBRL"`
	}
	var a result
	err = json.Unmarshal([]byte(s), &a)
	if err != nil {
		return 0, fmt.Errorf("usd: unmarshal: %w", err)
	}

	return StringToValue(fmt.Sprintf("%v", a.USDBRL.Price)), nil
}

func USDToBRL(usd float64) float64 {
	brl, err := USD()
	if err != nil {
		return 0
	}
	return usd * brl
}
