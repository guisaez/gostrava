package gostrava

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type StravaSubscription struct {
	CallbackURL  string
	VerifyToken  string
	
	*StravaClient
}

type Subscription struct {
	ID int64 `json:"id"`
}

func (ss *StravaSubscription) CreateSubscription(ctx context.Context, callbackURL string, ) (*Subscription, error) {

	apiUrl := baseURL + "/push_notifications"

	// Create the request body as url params, this will be handled internally
	params := url.Values{}
	params.Set("client_id", ss.ClientID)
	params.Set("client_secret", ss.ClientSecret)
	params.Set("callback_url", ss.CallbackURL)
	params.Set("verify_token", ss.VerifyToken)

	var subscription Subscription
	err := ss.put(context.Background(),"", "application/x-www-form-urlencoded", apiUrl, params, subscription)
	if err != nil {
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
