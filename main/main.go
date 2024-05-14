package main

import (
	"context"
	"encoding/json"

	"log"
	"net/http"
)

func main() {

	cors := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			w.Header().Add("Access-Control-Allow-Origin", "*")
			w.Header().Add("Vary", "Origin")
			w.Header().Add("Access-Control-Allow-Methods", "GET, POST, OPTIONS, HEAD, PUT, DELETE")
			w.Header().Add("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.Header().Add("Access-Control-Allow-Credentials", "true")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}

	OAuth := StravaOAuth{
		ClientID:     "125983",
		ClientSecret: "5913eb1ab24a593a640624c22cd27c0b3d10fdae",
		CallbackURL:  "http://localhost:8080/callback",
		Scopes:       []string{Scopes.Read, Scopes.ActivityRead},
	}

	router := http.NewServeMux()

	router.HandleFunc("GET /oauth", func(w http.ResponseWriter, r *http.Request) {

		http.Redirect(w, r, OAuth.AuthCodeURL("test", false), http.StatusFound)
	})

	onSuccess := func(tokens *StravaOAuthResp, w http.ResponseWriter, r *http.Request) {
		refresh, err := OAuth.Refresh(tokens.RefreshToken, nil)
		if err != nil {
			log.Println("RefreshError")
		}

		client := NewStravaClient(refresh.AccessToken, nil)

		activities, err := client.ListAthleteActivities(context.Background(), nil)
		if err != nil {
			panic(err)
		}

		for _, sa := range activities {
			ac, err := client.GetActivity(context.Background(), sa.ID, true)
			if err != nil {
				panic(err)
			}

			jsonD, err := json.MarshalIndent(ac, "", "\t")
			if err != nil {
				panic(err)
			}
			log.Println(string(jsonD))
			break;
		}



	}
	onFailure := func(err error, w http.ResponseWriter, r *http.Request) {
		// Handle error
		log.Println("Error Path")
		log.Println(err)
	}

	router.HandleFunc("GET /callback", OAuth.HandlerFunc(onSuccess, onFailure))

	

	srv := http.Server{
		Addr:    ":8080",
		Handler: cors(router),
	}

	srv.ListenAndServe()
}