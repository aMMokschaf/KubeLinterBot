//Package authentication is responsible for registering KubeLinterBot to a github-Repository.
//It also handles functions related to the oauth-token like serializing it or reading it again.
package authentication

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

//TODO rename this
type Client struct {
	jsonToken           string
	personalAccessToken string
	oauthClient         *http.Client
	GithubClient        *github.Client
}

//GetToken returns the access-token as a string, without bearer and expiry.
func (ao *Client) GetToken() string {
	return ao.personalAccessToken
}

//GetToken returns the access-token as a JSON-string, with bearer and expiry.
func (ao *Client) GetJSONToken() string {
	return ao.jsonToken
}

func (ao *Client) SetJSONToken(token string) {
	ao.jsonToken = token
}

func (ao *Client) getGithubClient() *github.Client {
	return ao.GithubClient
}

func (ao *Client) getOAuthClient() *http.Client {
	return ao.oauthClient
}

//TokenSource
type TokenSource struct {
	AccessToken string
}

//Token
func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

//CreateClient creates and returns the authObject needed to login to github.
func CreateClient(token string) *Client {
	var c Client
	c.personalAccessToken = extractTokenStringFromJSONToken(token)

	tokenSource := &TokenSource{
		AccessToken: c.personalAccessToken,
	}

	c.oauthClient = oauth2.NewClient(context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{ /*Transport: loggingRoundTripper{}*/ }), tokenSource)
	c.GithubClient = github.NewClient(c.oauthClient)

	return &c
}

type jsonTokenStruct struct {
	Access_token string
	Token_type   string
	Expiry       string
}

//extractTokenStringFromJSONToken extracts the access_token-String from the complete token.
func extractTokenStringFromJSONToken(completeToken string) string {
	var tokenStruct jsonTokenStruct
	json.Unmarshal([]byte(completeToken), &tokenStruct)
	var tokenString string = tokenStruct.Access_token
	return tokenString
}

// type loggingRoundTripper struct{}

// func (loggingRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
// 	fmt.Printf("Request: %s %s\n", req.Method, req.URL)
// 	if req.Body != nil {
// 		body, err := ioutil.ReadAll(req.Body)
// 		if err != nil {
// 			return nil, err
// 		}
// 		fmt.Printf("Body: %s\n", body)
// 		req.Body = ioutil.NopCloser(bytes.NewReader(body))
// 	}
// 	resp, err := http.DefaultTransport.RoundTrip(req)
// 	fmt.Printf("Response: %#v, %v\n", resp, err)
// 	return resp, err
// }
