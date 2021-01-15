package postcomment

import (
	"fmt"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var personalAccessToken string

type TokenSource struct {
	AccessToken string
}

//TODO: Doc and ask Malte wtf
func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

//PostComment TODO
func PostComment(token string, username string, reponame string, commitSha string, result []byte) {
	//fmt.Println("Entering PostComment")

	personalAccessToken = token
	tokenSource := &TokenSource{
		AccessToken: personalAccessToken,
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := github.NewClient(oauthClient)

	var bdy string = string(result)
	comment := github.RepositoryComment{Body: &bdy}
	//fmt.Println(username, reponame, commitSha, comment)
	_, r, err := client.Repositories.CreateComment(oauth2.NoContext, username, reponame, commitSha, &comment)
	if err != nil {
		fmt.Println("Posting kubelinters comment failed, error:", err)
		//return
	}
	fmt.Println(r)
}
