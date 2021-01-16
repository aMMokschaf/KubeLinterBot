//postcomment handles the posting of comments with KubeLinter's linting-results to the appropriate commit.
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

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

//PostComment authorizes with github to post KubeLinter's results to the commit.
func PostComment(token string, username string, reponame string, commitSha string, result []byte) {
	personalAccessToken = token
	tokenSource := &TokenSource{
		AccessToken: personalAccessToken,
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := github.NewClient(oauthClient)

	var bdy string = string(result)
	comment := github.RepositoryComment{Body: &bdy}
	_, r, err := client.Repositories.CreateComment(oauth2.NoContext, username, reponame, commitSha, &comment)
	if err != nil {
		fmt.Println("Posting kubelinters comment failed, error:", err)
	}
	fmt.Println(r)
}
