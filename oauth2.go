package gostrava

// import (
// 	"context"
// 	"fmt"
// 	"net"
// 	"net/http"
// 	"net/url"
// 	"strings"
// 	"time"
// )

// const OAuthBaseURL string = "https://www.strava.com/oauth/"

// const (
// 	oauthAuthorizationRequestSlug string = "authoriza"
// 	oauthGenerateTokenRequestSlug string = "token"
// 	oauthRevokeTokenRequestSlug   string = "revoke"
// )

// type ScopeString string

// const (
// 	Read            ScopeString = "read"
// 	ReadAll         ScopeString = "read_all"
// 	ProfileReadAll  ScopeString = "profile:read_all"
// 	ProfileWrite    ScopeString = "profile:write"
// 	ActivityRead    ScopeString = "activity:read"
// 	ActivityReadAll ScopeString = "activity:read_all"
// 	ActivityWrite   ScopeString = "activity:write"
// )

// type OAuth struct {
// 	clientID     string
// 	clientSecret string
// }

// // Retrieves the current clientID used by the OAuth2 client.
// func (c *Client) ClientID() string {
// 	return c.oauth.clientID
// }

// // Retrieves the clientSecret used by the OAuth2 client.
// func (c *Client) ClientSecret() string {
// 	return c.oauth.clientSecret
// }

// // Sets a new clientID for the OAuth2 client.
// func (c *Client) SetClientID(clientID string) {
// 	c.oauth.clientID = clientID
// }

// // Sets a new clientSecret for the OAuth2 client.
// func (c *Client) SetClientSecret(clientSecret string) {
// 	c.oauth.clientSecret = clientSecret
// }

// // AuthorizationCodeURL generates the authentication URL that the user will be redirected to
// // in order to initiate the OAuthFlow.
// //
// // Args:
// //   - redirectURI: The URI to which the authorization server will redirect the user-agent
// //     after the user grants or denies permission. This URI must be registered with the
// //     authorization server as part of the client registration.
// //   - state: An opaque value used by the client to maintain state between the request and callback.
// //     Typically used to prevent CSRF attacks and to maintain user state.
// //   - force: If true, it will force the authorization server to prompt the user for consent
// //     even if they have already done so for the current application.
// //   - scopes: A list of scopes (permissions) that the application requests access to.
// //
// // Returns:
// //   - string: The fully formed URL that the user-agent should be redirected to initiate the OAuthFlow.
// func (c *Client) AuthorizationCodeURL(redirectURI, state string, force bool, scopes []ScopeString) string {
// 	return MakeAuthorizationCodeURL(c.oauth.clientID, c.oauth.clientSecret, redirectURI, state, force, scopes)
// }

// type serverMessage struct {
// 	Type    string
// 	Message string
// 	Scopes  []ScopeString
// }

// func (c *Client) AuthorizationCodeRequest(redirectURI string, state string, force bool, scopes []ScopeString) (*AccessTokenResponse, error) {

// 	authURL := c.AuthorizationCodeURL(redirectURI, state, force, scopes)

// 	srvChan := make(chan int8)
// 	srvResponse := make(chan serverMessage, 1)
// 	errChan := make(chan serverMessage, 1)

// 	var srv *http.Server

// 	localRedirect := strings.Contains(redirectURI, "localhost") || strings.Contains(redirectURI, "127.0.0.1")
// 	if localRedirect {
// 		// Starts a localhost server that will handle the redirect url
// 		u, err := url.Parse(redirectURI)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to parse redirect URI: %s", err)
// 		}
// 		_, port, err := net.SplitHostPort(u.Host)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to split the redirect uri into host and port segments: %s", err)
// 		}

// 		srv := &http.Server{Addr: ":" + port}

// 		redirectPath := u.EscapedPath()
// 		http.HandleFunc(redirectPath, func(w http.ResponseWriter, r *http.Request) {

// 			q := r.URL.Query()

// 			errs := q.Get("error")
// 			if errs == "" {
// 				w.Write([]byte("Got an error from the server: " + errs))
// 				errChan <- serverMessage{Type: "error", Message: errs}
// 				return
// 			}

// 			code, scopes := q.Get("code"), strings.Split(q.Get("scope"), ",")

// 			w.Write([]byte("if you see this, code and scopes have been retrived! you can close this window"))

// 			var scopesString []ScopeString
// 			for _, s := range scopes {
// 				scopesString = append(scopesString, ScopeString(s))
// 			}

// 			srvResponse <- serverMessage{
// 				Type:    "code",
// 				Message: code,
// 				Scopes:  scopesString,
// 			}
// 		})

// 		go func() {
// 			srvChan <- 1
// 			err := srv.ListenAndServe()
// 			if err != nil && err != http.ErrServerClosed {
// 				fmt.Printf("error while serving localy: %s\n, err")
// 			}
// 		}()

// 		<-srvChan
// 	}

// 	fmt.Println("Redirect, go to the follwoing authorization URL to begin OAuth2 flow: \n %s\n\n", authURL)

// 	code, userAgreedScopes := "", []ScopeString{}

// 	var srvRes serverMessage

// 	if localRedirect {
// 		// we wait for the code to be returned by the server
// 		srvRes = <-srvResponse

// 		if srvRes.Type == "error" {
// 			return nil, fmt.Errorf("authroization error: %s\n", srvRes.Message)
// 		}

// 		code = srvRes.Message
// 		userAgreedScopes = srvRes.Scopes

// 		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 		defer cancel()
// 		if err := srv.Shutdown(ctx); err != nil {
// 			fmt.Printf("error while shutting down local server: %s\n", err)
// 		}
// 	} else {
// 		fmt.Println("Enter code and press enter:\n")
// 		_, err := fmt.Scan(&code)
// 		if err != nil {
// 			return fmt.Errorf("failed to read code from input: %s", err)
// 		}
// 	}

// 	if code == "" {
// 		return fmt.Errorf("No code was recevied from oAuth2 flow")
// 	}


// }

// type AccessTokenResponse struct {
// 	AccessToken  string `json:"access_token"`
// 	ExpiresAt    int64  `json:"expires_at"`    // The number of seconds since the epoch when the provided access token will expire
// 	ExpiresIn    int    `json:"expires_in"`    // Seconds until the short-lived access token will expire
// 	RefreshToken string `json:"refresh_token"` // The refresh token for this user, to be used to get the next access token for this user. Please expect that this value can change anytime you retrieve a new access token. Once a new refresh token code has been returned, the older code will no longer work
// 	TokenType    string `json:"token_type"`    // Bearer
// 	// Athlete      *gostrava.AthleteSummary `json:"athlete,omitempty"`    // A summary of the athlete information
// 	// Scopes       []Scope                  `json:"scopes,omitempty"`     // Scopes the user accepted
// }

// // MakeAuthorizationCodeURL generates the authentication URL that the user will be redirected to
// // in order to initiate the OAuthFlow.
// //
// // Args:
// //   - clientID: The ID assigned to the client application by the authorization server.
// //   - clientSecret: The client secret used for authentication with the authorization server.
// //   - redirectURI: The URI to which the authorization server will redirect the user-agent
// //     after the user grants or denies permission. This URI must be registered with the
// //     authorization server as part of the client registration.
// //   - state: An opaque value used by the client to maintain state between the request and callback.
// //     Typically used to prevent CSRF attacks and to maintain user state.
// //   - force: If true, it will force the authorization server to prompt the user for consent
// //     even if they have already done so for the current application.
// //   - scopes: A list of scopes (permissions) that the application requests access to.
// //
// // Returns:
// //   - string: The fully formed URL that the user-agent should be redirected to initiate the OAuthFlow.
// func MakeAuthorizationCodeURL(
// 	clientID, clientSecret, redirectURI, state string,
// 	force bool, scopes []ScopeString,
// ) string {
// 	q := url.Values{}
// 	q.Set("response_type", "code")
// 	q.Set("client_id", clientID)
// 	q.Set("redirect_uri", redirectURI)
// 	q.Set("scope", joinScopes(scopes))
// 	q.Set("state", state)

// 	if force {
// 		q.Set("approval_prompt", "force")
// 	}

// 	return fmt.Sprintf("%s%s?%s", OAuthBaseURL, oauthAuthorizationRequestSlug, q.Encode())
// }

// func joinScopes(scopes []ScopeString) string {
// 	stringScopes := make([]string, len(scopes))
// 	for i, scope := range scopes {
// 		stringScopes[i] = string(scope)
// 	}

// 	return strings.Join(stringScopes, ",")
// }
