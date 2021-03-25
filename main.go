package main

import (
	"browser/googleoauth2"
	"flag"
	"fmt"
	"net/http"
	"strings"
)

var client_id, cliet_secret, redirect_uri string
var port = 8080

func main() {
	flag.StringVar(&client_id, "id", "", "")
	flag.StringVar(&cliet_secret, "secret", "", "")
	flag.StringVar(&redirect_uri, "redirectUri", "", "")
	flag.Parse()
	fmt.Println("======")
	fmt.Printf("client_id    => %s\n", client_id)
	fmt.Printf("cliet_secret => %s\n", cliet_secret)
	fmt.Printf("redirect_uri => %s\n", redirect_uri)
	fmt.Println("======")

	o := googleoauth2.New(client_id, cliet_secret, redirect_uri)
	o.ApplyCode()

	http.Handle("/", &helloHandler{
		callback: o.ApplyToken,
	})
	go http.ListenAndServe(fmt.Sprintf(":%v", port), nil)

	fmt.Scanln()
}

type helloHandler struct {
	callback func(code string) (string, string)
}

func (h *helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var sub, idToken string
	if r.URL.Path == strings.ReplaceAll(redirect_uri, fmt.Sprintf("http://localhost:%v", port), "") {
		q := r.URL.Query()
		sub, idToken = h.callback(q["code"][0])
	}
	text := fmt.Sprintf("Sub => %s\n\nId Token => %s\n", sub, idToken)
	w.Write([]byte(text))
}
