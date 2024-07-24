package gostrava

import "net/http"

const BaseURLV3 = "https://www.strava.com/api/v3"

type Client struct {
	oauth OAuth

	baseURL string
	client *http.Client
}
