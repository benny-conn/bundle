package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"text/template"

	bundle "github.com/bennycio/bundle/internal"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("../../assets/templates/*.gohtml"))
}

func HandleRoot(w http.ResponseWriter, req *http.Request) {

	err := tpl.ExecuteTemplate(w, "index.gohtml", nil)
	if err != nil {
		log.Fatal(err)
	}

}
func HandleSignup(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {

		r.ParseForm()

		authUrl := "https://bundle.us.auth0.com/oauth/token"

		secret := os.Getenv("CLIENT_SECRET")

		values := url.Values{
			"grant_type":    {"client_credentials"},
			"client_id":     {"22oXY4A0h9Rfbo3XEAn8Fbptx715dBe4"},
			"client_secret": {secret},
			"audience":      {"https://bundlemc.io/auth/users"},
		}

		authRes, err := http.PostForm(authUrl, values)
		if err != nil {
			log.Fatal(err)
		}

		defer authRes.Body.Close()
		body, _ := ioutil.ReadAll(authRes.Body)

		auth := &bundle.Authorization{}

		err = json.Unmarshal(body, auth)
		if err != nil {
			log.Fatal(err)
		}

		newUser := bundle.User{
			r.FormValue("username"),
			r.FormValue("email"),
			r.FormValue("password"),
		}

		asJson, _ := json.Marshal(newUser)

		br := bytes.NewReader(asJson)

		createUser, err := http.NewRequest(http.MethodPost, "http://localhost:8070/users", br)
		if err != nil {
			log.Fatal(err)
		}

		createUser.Header.Set("authorization", "Bearer "+auth.Token)

		createRes, err := http.DefaultClient.Do(createUser)
		if err != nil {
			log.Fatal(err)
		}

		defer createRes.Body.Close()
		createResBody, _ := ioutil.ReadAll(createRes.Body)
		fmt.Println(string(createResBody))

	}

	err := tpl.ExecuteTemplate(w, "signup.gohtml", nil)
	if err != nil {
		log.Fatal(err)
	}

}
