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
	AuthURL      string `json:"auth_url"`
	TokenURL     string `json:"token_url"`
	UserURL      string `json:"user_url"`
}

func main() {

	fmt.Println("Welcome: this is a HTTP server project that serves as a client to test Keycloak auth code grant flow:")

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

	// Create a new redirect route route
	http.HandleFunc(cc.RedirectURI, func(w http.ResponseWriter, r *http.Request) {
		// First to get the value of the "code"
		err := r.ParseForm()
		if err != nil {
			fmt.Fprintf(os.Stdout, "could not parse query: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		code := r.FormValue("code")

		// Next, exchange code for token

		config := oauth2.Config{
			ClientID:     cc.ClientId,
			ClientSecret: cc.ClientSecret,
			RedirectURL:  cc.RedirectURL + cc.RedirectURI,
			Endpoint: oauth2.Endpoint{
				TokenURL: cc.TokenURL,
			},
		}
		ctx := context.Background()

		token, err := config.Exchange(ctx, code)
		if err != nil {
			http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
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
		fmt.Println("TONKEN for User INFO using: " + cc.UserURL)
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
