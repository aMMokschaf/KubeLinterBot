package getcommit

import (
	"fmt"
	"io"
	"net/http"
	"os"

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

	fmt.Println("Entering PostComment")

	personalAccessToken = token
	tokenSource := &TokenSource{
		AccessToken: personalAccessToken,
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := github.NewClient(oauthClient)

	/*
		var bdy string = string(result)
		comment := github.RepositoryComment{Body: &bdy}
	*/
	//fmt.Println(username, reponame, commitSha, comment)

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

	//authenticate with github
	//download whole commit
	return downloadStatus
}

func DownloadFile(url string, filename string) error {
	//TODO implement subfolders
	fmt.Println("Downloading file " + filename + "\n")
	const folder = "./downloadedYaml/"
	out, err := os.Create(folder + filename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url + filename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
