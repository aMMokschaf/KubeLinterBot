//Package authentication is responsible for registering KubeLinterBot to a github-Repository.
//It also handles functions related to the oauth-token like serializing it or reading it again.
package authentication

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"

	"main/internal/config"
)

var (
	// You must register the app at https://github.com/settings/applications
	// Set callback to http://127.0.0.1:7000/github_oauth_cb
	// Set ClientId and ClientSecret to
	oauthConf = &oauth2.Config{
		ClientID:     "c1aa344cbfeccc829f5e",
		ClientSecret: "146f4df2e4c12117156472d7620270281425b58b",
		// select level of access you want, refer to: https://developer.github.com/v3/oauth/#scopes
		Scopes:   []string{"user:email", "repo"},
		Endpoint: githuboauth.Endpoint,
	}
	// random string for oauth2 API calls to protect against CSRF
	oauthStateString = "thisshouldberandom" //TODO generate random string
)

//TODO struct
var s http.Server

//var waitGroup *sync.WaitGroup

const htmlIndex = `<html><body>
Logging in with <a href="/login">GitHub</a>
</body></html>
`

//RunAuth is called if KubeLinterBot is not authorized.
func RunAuth(cfg config.Config) { //wg *sync.WaitGroup) {
	//waitGroup = wg
	m := http.NewServeMux()
	s := &http.Server{Addr: ":7000", Handler: m}
	m.HandleFunc("/", handleMain)
	m.HandleFunc("/login", handleGitHubLogin)
	m.HandleFunc("/github_oauth_cb", handleGitHubCallback)
	m.HandleFunc("/shutdown", func(w http.ResponseWriter, req *http.Request) {
		handleShutdown(w, req) //, wg)
	})
	fmt.Print("Started running on http://127.0.0.1:7000\n")
	fmt.Println(s.ListenAndServe())
}

// /
func handleMain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlIndex))
}

// /login
func handleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	url := oauthConf.AuthCodeURL(oauthStateString, oauth2.AccessTypeOnline)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// /github_oauth_cb. Called by github after authorization is granted.
func handleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")
	token, err := oauthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		fmt.Printf("oauthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	oauthClient := oauthConf.Client(oauth2.NoContext, token)
	client := github.NewClient(oauthClient)
	user, _, err := client.Users.Get(oauth2.NoContext, "")
	if err != nil {
		fmt.Printf("client.Users.Get() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	jsonToken, err := tokenToJSON(token)
	if err != nil {
		fmt.Println("Problem with serializing token to JSON:", err)
	}

	cfg, err := config.OptionParser()
	if err != nil {
		panic(err)
	}
	cfg.User.AccessToken = jsonToken
	config.WriteOptionsToFile(*cfg)

	//fmt.Println(token, jsonToken)

	fmt.Printf("Logged in as GitHub user: %s\n", *user.Login)
	//fmt.Println("TOKEN:", jsonToken)
	http.Redirect(w, r, "/shutdown", http.StatusTemporaryRedirect)
}

//handleShutdown is the handler for /shutdown.
func handleShutdown(w http.ResponseWriter, r *http.Request) { //, wg *sync.WaitGroup) {
	//wg.Done()
	s.Shutdown(context.Background())
}

//tokenToJSON converts a oauth2.Token to a JSON-String.
func tokenToJSON(token *oauth2.Token) (string, error) {
	if d, err := json.Marshal(token); err != nil {
		return "", errors.Wrap(err, "marshaling token as JSON")
	} else {
		return string(d), nil
	}
}

//tokenFromJSON parses a JSON-string to a oauth2.Token.
func tokenFromJSON(jsonStr string) (*oauth2.Token, error) {
	var token oauth2.Token
	if err := json.Unmarshal([]byte(jsonStr), &token); err != nil {
		return nil, err
	}
	return &token, nil
}
