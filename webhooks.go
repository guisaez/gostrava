package gostrava

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const (
	subscriptionEndpoint string = "push_subscription"
)

// SubscriptionService provides methods for creating, viewing and deleting webhook subscription.
type SubscriptionService struct {
	service
}

type Subscription struct {
	ID int `json:"id"` // Subscription ID
}

func (s *SubscriptionService) Subscribe(ctx context.Context, callbackURL, verifyToken string) (*Subscription, *http.Response, error) {
	subscriptionURL := s.SubscriptionRequestURL(callbackURL, verifyToken)

	req, err := s.client.NewRequest(http.MethodPost, subscriptionURL, nil, func(req *http.Request) error {
		req.URL, _ = url.Parse(subscriptionURL)
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	subscription := new(Subscription)
	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return nil, resp, err
	}

	return subscription, resp, nil
}

func (s *SubscriptionService) Delete(ctx context.Context, id int) (bool, *http.Response, error) {
	q := url.Values{}
	q.Set("client_id", s.client.clientID)
	q.Set("client_secret", s.client.clientSecret)

	urlStr := fmt.Sprintf("%s/%d?%s", subscriptionEndpoint, id, q.Encode())

	req, err := s.client.NewRequest(http.MethodDelete, urlStr, nil, nil)
	if err != nil {
		return false, nil, err
	}

	resp, err := s.client.Do(ctx, req)
	if err != nil {
		return false, resp, err
	}

	return resp.StatusCode == http.StatusNoContent, resp, err
}

type subscriptionValidationResponse struct {
	Challenge string `json:"hub.challenge"` // Random string the callback address must echo back to verify its existence.
}

func SubscriptionValidationHandler(verifyToken string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		mode := q.Get("hub.mode")
		challenge := q.Get("hub.challenge")
		incomingVerifyToken := q.Get("hub.verify_token")

		if verifyToken != incomingVerifyToken {
			return
		}

		if mode != "subscribe" {
			return
		}

		json, _ := json.Marshal(subscriptionValidationResponse{Challenge: challenge})

		w.WriteHeader(http.StatusOK)
		w.Write(json)
	}
}

func (s *SubscriptionService) SubscriptionRequestURL(callbackURL, verifyToken string) string {
	return BuildSubscriptionRequestURL(s.client.BaseURL.String(), s.client.clientID, s.client.clientSecret, callbackURL, verifyToken)
}

// BuildSubscriptionCreate URL constructs the URL for creating a webhook subscription.
func BuildSubscriptionRequestURL(baseURL, clientID, clientSecret, callbackURL, verifyToken string) string {
	q := url.Values{}
	q.Set("client_id", clientID)
	q.Set("client_secret", clientSecret)
	q.Set("callback_url", callbackURL)
	q.Set("verify_token", verifyToken)

	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	return fmt.Sprintf("%s%s?%s", baseURL, subscriptionEndpoint, q.Encode())
}
