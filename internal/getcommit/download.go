//getcommit is used to download all folders with .yaml and .yml-files.
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

const mainFolder = "./downloadedYaml/"

type TokenSource struct {
	AccessToken string
}

//Token creates the oauth2.Token for oauth.
func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

//DownloadCommit authenticates with oauth and downloads all folders with .yaml or .yml-files.
//These are then passed to the KubeLinter-binary.
func DownloadCommit(token string, username string, reponame string, commitSha string, addedFiles []string, modifiedFiles []string) bool {
	var downloadStatus = false

	personalAccessToken = token
	tokenSource := &TokenSource{
		AccessToken: personalAccessToken,
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := github.NewClient(oauthClient)

	folder, err := downloadFolder("", username, reponame, client, oauthClient)
	if err != nil {
		fmt.Println("Error while creating folder.", err)
	} else {
		fmt.Println(folder)
	}

	downloadStatus = true
	return downloadStatus
}

//downloadFolder downloads all files in a folder, creating subfolders as necessary.
func downloadFolder(path string, username string, reponame string, client *github.Client, oauthClient *http.Client) ([]*github.RepositoryContent, error) {
	var options = github.RepositoryContentGetOptions{}
	file, folder, r, err := client.Repositories.GetContents(oauth2.NoContext,
		username,
		reponame,
		path,
		&options)
	if err != nil {
		return folder, err
	} else {
		fmt.Println("\nfolder:", folder, "\nfile:", file, "\nresponse:", r)
		//	return folder, nil
	}
	for _, file := range folder {
		if string(file.GetType()) == "dir" {
			err := os.MkdirAll(string(mainFolder+file.GetPath()), 0755)
			if err != nil {
				fmt.Println("Error while creating folder.", err)
				//return downloadStatus
			} else {
				fmt.Println("Folder created", file.GetPath())
			}
			downloadFolder(file.GetPath(), username, reponame, client, oauthClient)
		} else if string(file.GetType()) == "file" {
			downloadFile(file.GetDownloadURL(), file.GetPath())
		}
	}
	return folder, err
}

//downloadFile downloads a single file.
func downloadFile(url string, filepath string) error {
	fmt.Println("Downloading file " + url + "\n")
	out, err := os.Create(mainFolder + filepath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
