# Go Strava Client Wrapper & Authentication Package

This Go package provides a convenient wrapper for interacting with the Strava API, along with authentication functionalities to streamline the integration of Strava's services into your Go applications.

Features

* Simple Integration: Easily incorporate Strava's API into your Go applications with intuitive wrappers for common API endpoints.

* Authentication: Securely authenticate users with Strava OAuth 2.0, handling token generation and refreshing seamlessly.
    
* Efficient Data Retrieval: Retrieve activity data, athlete information, and more from Strava's API with minimal overhead.
    
* Flexible Configuration: Configure the package according to your needs, including setting authentication parameters and API endpoints.

## Import 

```Shell
go get github.com/guisaez/go-strava
```

## Usage

### Client

```go
func main() {

    // Application Credentials can be obtained here https://www.strava.com/settings/api
    clientId, clientSecret := <clientID>, <clientSecret>

    // The third argument corresponds to a custom http.Client, if nil it will use the default http.Client
    strava := NewStravaClient(clientId, clientSecret, nil)

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
func main(){

    stravaClient := NewStravaClient(clientId, clientSecret, nil)

    athletes = &StravaAthletes{
        AccessToken: <accessToken>,
        StravaClient: stravaClient
    }

    athlete, err := athletes.CurrentAthlete(context.Background())
    ...
}
```