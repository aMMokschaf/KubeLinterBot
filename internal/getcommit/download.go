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

//DownloadCommit authenticates with oauth and TODO
func DownloadCommit(token string, username string, reponame string, commitSha string, addedFiles []string, modifiedFiles []string) bool {
	var downloadStatus = false

	//fmt.Println("Entering DownloadCommit")

	personalAccessToken = token
	tokenSource := &TokenSource{
		AccessToken: personalAccessToken,
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := github.NewClient(oauthClient)

	var options = github.RepositoryContentGetOptions{}
	file, folder, r, err := client.Repositories.GetContents(oauth2.NoContext,
		username,
		reponame,
		"",
		&options)
	if err != nil {
		fmt.Println("GetCommit failed, error:", err)
		return downloadStatus
	}
	fmt.Println("\nfolder:", folder, "\nfile:", file, "\nresponse:", r)

	//TODO folders cant be downloaded yet
	for _, filename := range modifiedFiles {
		for _, filename2 := range folder {
			fmt.Println(filename)
			//folder ist ein array
			//jedes modifiedFile mit allen paths abgleichen?
			//Dann ggf an DownloadFile übergeben
			downloadFile2(filename2.GetDownloadURL(), filename2.GetName())
		}
	}

	downloadStatus = true

	return downloadStatus
}

func downloadFile2(url string, filename string) error {
	//TODO implement subfolders
	fmt.Println("Downloading file " + url + "\n")
	const mainFolder = "./downloadedYaml/"
	out, err := os.Create(mainFolder + filename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
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

func downloadFile(url string, filename string) error {
	//TODO implement subfolders
	fmt.Println("Downloading file " + filename + "\n")
	const mainFolder = "./downloadedYaml/"
	out, err := os.Create(mainFolder + filename)
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
