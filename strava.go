package gostrava

import "net/http"

const baseURLv3 = "https://www.strava.com/api/v3"

type Client struct {
	BaseURL string

	client *stravaHTTP
}