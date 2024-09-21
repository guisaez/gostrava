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
	subscriptionEndpoint string = "/api/v3/push_subscriptions"
)

// SubscriptionService provides methods for creating, viewing and deleting webhook subscription.
type SubscriptionService service

type Subscription struct {
	ID int `json:"id"` // Subscription ID
}

func (s *SubscriptionService) Subscribe(ctx context.Context, callbackURL, verifyToken string) (*Subscription, *http.Response, error) {
	subscriptionURL := s.SubscriptionRequestURL(callbackURL, verifyToken)

	fmt.Println(subscriptionURL)
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

// SubscriptionValidationHandler validates the subscription request using a token
// provided by the caller and responds with the appropriate challenge.
func SubscriptionValidationHandler(verifyToken string) http.HandlerFunc {
	// Return an HTTP handler function
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract query parameters from the URL
		q := r.URL.Query()

		// Get the "hub.mode", "hub.challenge", and "hub.verify_token" from the query parameters
		mode := q.Get("hub.mode")
		challenge := q.Get("hub.challenge")
		incomingVerifyToken := q.Get("hub.verify_token")

		// Validate that the provided verify_token matches the expected one
		if verifyToken != incomingVerifyToken {
			// If tokens don't match, return a 401 Unauthorized response
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Check if the mode is "subscribe" as expected for a subscription request
		if mode != "subscribe" {
			// If the mode is not "subscribe", return a 400 Bad Request response
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// If the verification and mode checks pass, respond with the challenge token
		response := subscriptionValidationResponse{Challenge: challenge}

		// Write a 200 OK status to indicate success
		w.WriteHeader(http.StatusOK)

		// Send the challenge response as JSON back to the client
		json.NewEncoder(w).Encode(response)
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
