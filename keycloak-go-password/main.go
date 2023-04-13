package main

import (
	"fmt"
	"net/http"
	"os"

	"context"
	"encoding/json"

	"golang.org/x/oauth2"
)

type ClientConfig struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectURL  string `json:"redirect_url"`
	RedirectURI  string `json:"redirect_uri"`
	LoginURI     string `json:"login_uri"`
	AuthURL      string `json:"auth_url"`
	TokenURL     string `json:"token_url"`
	UserURL      string `json:"user_url"`
}

func main() {
	fmt.Println("Welcome: this is a HTTP server project that serves as a client to test Keycloak password grant flow:")

	cf, _ := os.ReadFile("public/client.json")
	cc := ClientConfig{}
	if err := json.Unmarshal(cf, &cc); err != nil {
		panic(err)
	}
	fmt.Printf("Keycloak client ID: %s.\n", cc.ClientId)
	fmt.Printf("Keycloak token URL: %s.\n", cc.TokenURL)
	fmt.Printf("Keycloak user information URL: %s.\n", cc.UserURL)
	fmt.Printf("Default home URL: %s.\n", "http://localhost:8082/")

	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	config := oauth2.Config{
		ClientID:     cc.ClientId,
		ClientSecret: cc.ClientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL: cc.TokenURL,
		},
	}

	// A login request will be submitted to below action
	http.HandleFunc("/oauth/login", func(w http.ResponseWriter, r *http.Request) {
		// First to get the values of "user / pass"
		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not parse query: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		user := r.FormValue("username")
		pass := r.FormValue("password")

		// Next, exchange u/p for token

		ctx := context.Background()

		fmt.Println("username =" + user)

		pmx := ""

		if pass == "" || len(pass) <= 2 {
			pmx = "**"

		} else {

			pmx = pass[:2]

			for i := 0; i < len(pass)-2; i++ {
				pmx = pmx + "*"
			}
		}

		fmt.Printf("password = %s\n", pmx)

		token, err := config.PasswordCredentialsToken(ctx, user, pass)
		if err != nil {
			http.Error(w, "Failed to get token: "+err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("Access token: " + token.AccessToken)

		w.Header().Set("Location", "/welcome.html?access_token="+token.AccessToken)

		w.WriteHeader(http.StatusFound)
	})

	http.HandleFunc("/oauth/userinfo", func(w http.ResponseWriter, r *http.Request) {
		httpClient := http.Client{}
		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not parse query: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		token := r.FormValue("access_token")

		req, err := http.NewRequest(http.MethodPost, cc.UserURL, nil)
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not create HTTP request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// We set this header since we want the response
		// as JSON
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Authorization", "Bearer "+token)
		// Send out the HTTP request
		res, err := httpClient.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not send HTTP request: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer res.Body.Close()
		var t OAuthUserinfoResponse
		if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
			fmt.Fprintf(os.Stdout, "could not parse JSON response: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Println(t.Name)
		rep, _ := json.Marshal(&t)
		w.Write([]byte([]byte(rep)))

	})

	http.ListenAndServe(":8082", nil)
}

type OAuthUserinfoResponse struct {
	Name string `json:"name"`
}
