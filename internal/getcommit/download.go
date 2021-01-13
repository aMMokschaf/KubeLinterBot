package getcommit

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

func DownloadCommit(token string, username string, reponame string, commitSha string) bool {
	var downloadStatus = false

	fmt.Println("Entering DownloadCommit")

	personalAccessToken = token
	tokenSource := &TokenSource{
		AccessToken: personalAccessToken,
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := github.NewClient(oauthClient)

	var options = github.RepositoryContentGetOptions{}
	file, folder, r, err := client.Repositories.GetContents(oauth2.NoContext,
		"aMMokschaf",
		"yamls",
		"",
		&options)
	if err != nil {
		fmt.Println("GetCommit failed, error:", err)
		return downloadStatus
	}
	fmt.Println(folder, file, r)
	downloadStatus = true

	return downloadStatus
}
