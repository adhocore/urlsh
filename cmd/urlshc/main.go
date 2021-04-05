package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	ureq "github.com/adhocore/urlsh/request"
	usrv "github.com/adhocore/urlsh/service/url"
)

type URLOutput struct {
	Status     uint    `json:"status" binding:"required"`
	ShortCode  string  `json:"short_code"`
	Message    string  `json:"message"`
}

const path = "/api/urls"

func main() {
	// Set sensible defaults for client
	os.Setenv("APP_ALLOW_DUPE_URL", "1")
	os.Setenv("APP_CHECK_URL_REACH", "1")

	in, err := input()
	abort(err)

	out := request(in)
	info(in, out)
}

func input() (ureq.URLInput, error) {
	var i ureq.URLInput
	var words string

	flag.StringVar(&i.URL, "url", "", "The long URL")
	flag.StringVar(&i.ExpiresOn, "expires", "", "Expiry date (optional, format: yyyy-mm-dd)")
	flag.StringVar(&words, "keywords", "", "CSV keywords (optional, format: word-1,word2,word_3,...)")
	flag.Parse()

	if words != "" {
		words = strings.ReplaceAll(words, ", ", ",")
		words = strings.ReplaceAll(words, " ,", ",")
		i.Keywords = strings.Split(words, ",")
	}
	if i.ExpiresOn != "" && len(i.ExpiresOn) == 10 {
		i.ExpiresOn = i.ExpiresOn + " 23:59:59"
	}

	_, err := usrv.ValidateURLInput(i)

	return i, err
}

func abort(err error) {
	if err != nil {
		log.Fatalf("\033[31merror: %v\033[m\n", err)
	}
}

func request(in ureq.URLInput) URLOutput {
	buf, err := json.Marshal(in)
	abort(err)

	req, err := http.NewRequest("POST", host() + path, bytes.NewBuffer(buf))
	abort(err)

	client := &http.Client{}
	res, err := client.Do(req)
	abort(err)

	defer res.Body.Close()

	var o URLOutput
	err = json.NewDecoder(res.Body).Decode(&o)
	abort(err)

	return o
}

func host() string {
	host := os.Getenv("URLSH_HOST")
	if host == "" {
		return "https://urlssh.xyz"
	}

	if strings.Index(host, "http") != 0 {
		host = "http://" + host
	}

	return host
}

func info(i ureq.URLInput, o URLOutput) {
	if o.ShortCode == "" {
		log.Fatalf("\033[31merror: %v\033[m\n", o.Message)
	}

	fmt.Printf("\033[36mURL:\033[m \033[33m%s\033[m\n", i.URL)
	fmt.Printf("\033[36mShort URL:\033[m \033[32m%s\033[m\n", host() + "/" + o.ShortCode)
}
