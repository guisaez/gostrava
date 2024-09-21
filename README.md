# Strava V3 API Client v2

![Go Version](https://img.shields.io/badge/go-%3E=1.16-blue)

## Overview 

`gostrava` is a Go client library for interacting with "Strava" API. This client wraps all the API endpoints referenced in the [Official Strava API documentation](https://developers.strava.com/docs/reference/), providing an easy-to-sue interface for Go developers to integrate.

## What's different in V2

- Allows more flexibility for the developer. V1 was an experimental implementation of the service as I was diving into Golang.
- Each service is part of a single client, which allows to avoid re-declaring it just to handle authentication, etc. This allows to initialize alongside the server and make use of it as needed.
- Ability to customize requests using context.

## Adding `gostrava` to your project.
```
go get github.com/guisaez/gostrava/v2
```

## Authorization/Authentication

All request to the Strava API require authentication. A `mandatory` access token will need to be provided. 

Authorization tokens, will not be stored in the client, they will need to be provided for each function call.

### Documentation TO-DO