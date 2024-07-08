# Go Strava Client Wrapper & Authentication Package

This Go package provides a convenient wrapper for interacting with the Strava API, along with authentication functionalities to streamline the integration of Strava's services into your Go applications.

Features

* Simple Integration: Easily incorporate Strava's API into your Go applications with intuitive wrappers for common API endpoints.

* Authentication: Securely authenticate users with Strava OAuth 2.0, handling token generation and refreshing seamlessly.
    
* Efficient Data Retrieval: Retrieve activity data, athlete information, and more from Strava's API with minimal overhead.
    
* Flexible Configuration: Configure the package according to your needs, including setting authentication parameters and API endpoints.

## Import 

```Shell
go get github.com/guisaez/gostrava
```

## Usage

### Client

```go
import (
    gostrava "github.com/guisaez/gostrava"
)

func main() {

    // Application Credentials can be obtained here https://www.strava.com/settings/api
    clientId, clientSecret := <clientID>, <clientSecret>

    // The third argument corresponds to a custom http.Client, if nil it will use the default http.Client
   client := gostrava.NewClient(nil)

    ...
}
```

### Modules

APIFunctions are separated into different modules:

* Activities
* Athletes
* Clubs
* Gears
* Routes
* SegmentEfforts
* Segments
* Streams
* Uploads

Each module has its own set of allowed methods. Example:

```go
    access_token := "<access_token>"

    athlete, err := client.Athletes.GetAuthenticatedAthlete(access_token)

    athleteActivities, err := client.Activities.ListAthleteActivities(access_token)
```

### OAuth

#### Scopes
```go

    var StravaScopes = struct {
        Read            string
        ReadAll         string
        ProfileReadAll  string
        ProfileWrite    string
        ActivityRead    string
        ActivityReadAll string
        ActivityWrite   string
    }{
        "read",
        "read_all",
        "profile:read_all",
        "profile:write",
        "activity:read",
        "activity:read_all",
        "activity:write",
    }
```

```go

func main() {

    oauthOpts := gostrava.OAuthOpts{
        ClientID: "<client_id>"
        ClientSecret: "<client_secret>"
        CallbackURL: "http://localhost:8080/callback"
        Scopes: []string{gostrava.StravaScopes.Read, gostrava.StravaScopes.ActivityWrite}
    }

    oauth := gostrava.NewStravaOAuth(oauthOpts)

    oauth.MakeAuthCodeUrl(false, "")
}
```

Applying Client Redirection to Init OAuth Flow - Example

```go
func(w http.ResponseWriter, *http.Request) {
    ...

    // Redirects the user to the Strava Authorization Page
    w.Redirect(oauth.AuthCodeURL(true, ""), http.StatusFound)
}
```

Requesting Access and RefreshTokens

```go

    authCode := "123123"
    tokens, err := oauth.Exchange(code)
```

Refreshing the tokens
```go

    refreshToken := "a1231"

    tokens, err := oauth.Refresh(refreshToken)
```

Revoking Access

```go
    accessToken = "123421a"

    err := oauth.RevokeAccess(accessToken)
```

Storing Session Info

```go
func main() {

    onSuccess := func(tokens *StravaOAuthResponse, w http.ResponseWriter, r *http.Request) {
        // Add login here

        // Save tokens on Database etc
    }

    onFailure := func(err error, w http.ResponseWrite, r *http.Request) {
        // Handle the error here
    }

    router.HandleFunc("GET /callback", oauth.HandlerFunc(onSuccess, onFailure))
    
}
```