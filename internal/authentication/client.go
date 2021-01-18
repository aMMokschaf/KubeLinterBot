//Package authentication is responsible for registering KubeLinterBot to a github-Repository.
//It also handles functions related to the oauth-token like serializing it or reading it again.
package authentication

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var jsonToken string
var personalAccessToken string
var oauthClient *http.Client
var githubClient *github.Client

type TokenSource struct {
	AccessToken string
}

//GetToken returns the access-token as JSON.
func GetToken() string {
	return personalAccessToken
}

func GetFullToken() string {
	return jsonToken
}

type jsonTokenStruct struct {
	Access_token string
	Token_type   string
	Expiry       string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

//ExtractTokenStringFromJSONToken extracts the access_token-String from the complete token.
func extractTokenStringFromJSONToken(completeToken string) string {
	var tokenStruct jsonTokenStruct
	json.Unmarshal([]byte(completeToken), &tokenStruct)
	var tokenString string = tokenStruct.Access_token
	return tokenString
}

func CreateClient(token string) {
	personalAccessToken = extractTokenStringFromJSONToken(token)

	tokenSource := &TokenSource{
		AccessToken: personalAccessToken,
	}

	oauthClient = oauth2.NewClient(oauth2.NoContext, tokenSource)
	githubClient = github.NewClient(oauthClient)
	//Remove this if possible
	fmt.Println(githubClient)
}

func GetGithubClient() *github.Client {
	return githubClient
}

func GetOAuthClient() *http.Client {
	return oauthClient
}
