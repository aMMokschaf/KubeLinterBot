package authentication

//TODO: Doc the whole thing
import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

var (
	// You must register the app at https://github.com/settings/applications
	// Set callback to http://127.0.0.1:7000/github_oauth_cb
	// Set ClientId and ClientSecret to
	oauthConf = &oauth2.Config{
		ClientID:     "c1aa344cbfeccc829f5e",
		ClientSecret: "146f4df2e4c12117156472d7620270281425b58b",
		// select level of access you want https://developer.github.com/v3/oauth/#scopes
		Scopes:   []string{"user:email", "repo"},
		Endpoint: githuboauth.Endpoint,
	}
	// random string for oauth2 API calls to protect against CSRF
	oauthStateString = "thisshouldberandom"
)

var jsonToken string

var s http.Server
var waitGroup *sync.WaitGroup

const htmlIndex = `<html><body>
Logged in with <a href="/login">GitHub</a>
</body></html>
`

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

// /github_oauth_cb. Called by github after authorization is granted
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
	jsonToken, err = tokenToJSON(token)
	if err != nil {
		fmt.Println("Problem with serializing token to JSON:", err)
	}
	fmt.Printf("Logged in as GitHub user: %s\n", *user.Login)
	//fmt.Println("TOKEN:", jsonToken)
	http.Redirect(w, r, "/shutdown", http.StatusTemporaryRedirect)
}

func handleShutdown(w http.ResponseWriter, r *http.Request) {
	waitGroup.Done()
	s.Shutdown(context.Background())
}

func RunAuth(wg *sync.WaitGroup) {
	waitGroup = wg
	m := http.NewServeMux()
	s := &http.Server{Addr: ":7000", Handler: m}
	m.HandleFunc("/", handleMain)
	m.HandleFunc("/login", handleGitHubLogin)
	m.HandleFunc("/github_oauth_cb", handleGitHubCallback)
	m.HandleFunc("/shutdown", handleShutdown)
	fmt.Print("Started running on http://127.0.0.1:7000\n")
	fmt.Println(s.ListenAndServe())
}

func tokenToJSON(token *oauth2.Token) (string, error) {
	if d, err := json.Marshal(token); err != nil {
		return "", err
	} else {
		return string(d), nil
	}
}

func GetToken() string {
	return jsonToken
}

type jsonTokenStruct struct {
	Access_token string
	Token_type   string
	Expiry       string
}

//This method extracts the access_token-String from the complete token.
func ExtractTokenStringFromJSONToken(completeToken string) string {
	var tokenStruct jsonTokenStruct
	json.Unmarshal([]byte(completeToken), &tokenStruct)
	var tokenString string = tokenStruct.Access_token
	return tokenString
}

func tokenFromJSON(jsonStr string) (*oauth2.Token, error) {
	var token oauth2.Token
	if err := json.Unmarshal([]byte(jsonStr), &token); err != nil {
		return nil, err
	}
	return &token, nil
}
