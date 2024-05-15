package go_strava

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
)

type StravaSubscription struct {
	ClientID     string
	ClientSecret string
	CallbackURL  string
	VerifyToken  string
	// Generates an HTTP Client for making requests during the OAuth token exchange process.
	RequestClient func(r *http.Request) *http.Client
}

type Subscription struct {
	ID int64 `json:"id"`
}

func (ss *StravaSubscription) CreateSubscription(ctx context.Context, callbackURL string, ) (*Subscription, error) {

	apiUrl := baseURL + "/push_notifications"

	// Create the request body
	params := url.Values{}
	params.Set("client_id", ss.ClientID)
	params.Set("client_secret", ss.ClientSecret)
	params.Set("callback_url", ss.CallbackURL)
	params.Set("verify_token", ss.VerifyToken)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiUrl, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	var subscription Subscription
	if err := json.NewDecoder(resp.Body).Decode(&subscription); err != nil {
		return nil, err
	}

	return &subscription, nil
}

type CallbackResponse struct {
	Challenge string `json:"hub.challenge"`
}
 
func (ss *StravaSubscription) SubscriptionCallbackValidation(w http.ResponseWriter, r *http.Request) {

	query_params := r.URL.Query()

	if query_params.Get("hub.mode") != "subscribe" || query_params.Get("hub.verify_token") != ss.VerifyToken {
		return 
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	var payload = &CallbackResponse{
		Challenge: query_params.Get("hub.challenge"),
	}
	if err := json.NewEncoder(w).Encode(&payload); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
