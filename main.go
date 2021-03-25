package main

import (
	"flag"
	"fmt"
	"google-oauth2-bot/googleoauth2"
	"net/http"
	"os"
	"strings"
)

var client_id, cliet_secret, redirect_uri string
var port = 8080

func main() {
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "[usage]")
		fmt.Fprintln(os.Stderr, "-id=<value>          : Google oauth2 client_id")
		fmt.Fprintln(os.Stderr, "-secret=<value>      : Google oauth2 cliet_secret")
		fmt.Fprintln(os.Stderr, "-redirectUri=<value> : Google oauth2 redirect_uri")

		fmt.Fprintln(os.Stderr, "\nRef this URL to get correlation variables, https://developers.google.com/identity/protocols/oauth2")
	}
	flag.StringVar(&client_id, "id", "", "client_id")
	flag.StringVar(&cliet_secret, "secret", "", "cliet_secret")
	flag.StringVar(&redirect_uri, "redirectUri", "", "redirect_uri")
	flag.Parse()

	if len(client_id) == 0 || len(cliet_secret) == 0 || len(redirect_uri) == 0 {
		flag.Usage()
		return
	}

	fmt.Println("======")
	fmt.Printf("client_id    => %s\n", client_id)
	fmt.Printf("cliet_secret => %s\n", cliet_secret)
	fmt.Printf("redirect_uri => %s\n", redirect_uri)
	fmt.Println("======")
	//
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
