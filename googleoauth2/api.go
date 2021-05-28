package googleoauth2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
)

type GOAUTH2 struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

func New(client_id, client_secret, redirect_url string) *GOAUTH2 {
	return &GOAUTH2{
		ClientID:     client_id,
		ClientSecret: client_secret,
		RedirectURL:  redirect_url,
	}
}

/*
ApplyCode - 申請 code
*/
func (o *GOAUTH2) ApplyCode() {
	urlStr := "https://accounts.google.com/o/oauth2/v2/auth"
	scope := "email"
	access_type := "online"
	include_granted_scopes := "true"
	state := "state_parameter_passthrough_value"
	response_type := "code"
	redirect_uri := o.RedirectURL
	client_id := o.ClientID
	api := fmt.Sprintf("%s?scope=%s&access_type=%s&include_granted_scopes=%s&state=%s&response_type=%s&redirect_uri=%s&client_id=%s", urlStr, scope, access_type, include_granted_scopes, state, response_type, redirect_uri, client_id)
	//
	openbrowser(api)
}

func (o *GOAUTH2) ApplyToken(code string) (sub, idToken string) {
	idToken = o.GetIdToken(code)
	sub = o.GetSub(idToken)
	fmt.Printf("sub =>\n%s\n\n", sub)
	fmt.Printf("id token =>\n%s\n", idToken)
	return sub, idToken
}

func (o *GOAUTH2) GetIdToken(code string) string {
	/*
		curl -X POST \
		https://www.googleapis.com/oauth2/v4/token \
		  -H 'Content-Type: application/x-www-form-urlencoded' \
		  -d 'code=&client_id=&client_secret=&redirect_uri=&grant_type=authorization_code'
	*/

	urlStr := "https://www.googleapis.com/oauth2/v4/token"
	f := url.Values{}
	f.Add("code", code)
	f.Add("client_id", o.ClientID)
	f.Add("client_secret", o.ClientSecret)
	f.Add("redirect_uri", o.RedirectURL)
	f.Add("grant_type", "authorization_code")

	client := &http.Client{}
	res, err := client.PostForm(urlStr, f)
	if err != nil {
		panic(err)
	}
	//
	type RESP struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   string `json:"expires_in"`
		Scope       string `json:"scope"`
		TokenType   string `json:"token_type"`
		IdToken     string `json:"id_token"`
	}

	obj, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	//
	var r RESP
	json.Unmarshal(obj, &r)
	return r.IdToken
}

func (o *GOAUTH2) GetSub(id_token string) string {
	api := "https://oauth2.googleapis.com/tokeninfo?id_token=" + id_token
	res, err := http.Get(api)
	if err != nil {
		log.Fatal(err)
	}
	obj, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	//
	type RESP struct {
		Sub string `json:"sub"`
	}
	var r RESP
	json.Unmarshal(obj, &r)
	return r.Sub
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")

	}
	if err != nil {
		log.Fatal(err)
	}
}
