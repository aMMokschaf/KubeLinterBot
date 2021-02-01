//Package getcommit is used to download all folders with .yaml and .yml-files.
package getcommit

import (
	"context"
	"fmt"
	"io"
	"main/internal/authentication"
	"net/http"
	"os"

	"github.com/google/go-github/github"
)

const mainDir = "./downloadedYaml/"

//DownloadCommit authenticates with oauth and downloads all folders with .yaml or .yml-files.
//These are then passed to the KubeLinter-binary.
func DownloadCommit(ownername string, reponame string, commitSha string, branch string, filenames []string, number int, client authentication.Client) ([]*github.RepositoryContent, error) {
	repoContent, err := downloadFolder(ownername, reponame, "", commitSha, branch, client.GithubClient)
	if err != nil {
		return nil, err
	}
	return repoContent, err
}

//downloadFolder downloads all files in a folder, creating subfolders as necessary.
func downloadFolder(ownername string, reponame string, subpath string, commitSha string, branchRef string, client *github.Client) ([]*github.RepositoryContent, error) {
	fmt.Println("downloadFolder")
	var commitDir string
	if commitSha != "" {
		commitDir = mainDir + commitSha + "/"
	} else {
		commitDir = mainDir + ownername + reponame + "/"
	}

	branch, _, err := client.Repositories.GetBranch(context.Background(),
		ownername,
		reponame,
		branchRef)
	if err != nil {
		return nil, err
	}

	var options = github.RepositoryContentGetOptions{Ref: branch.GetName()}
	//f, err := client.Repositories.DownloadContents(context.Background(), ownername, reponame, subpath, &options)

	_, folder, _, err := client.Repositories.GetContents(context.Background(),
		ownername,
		reponame,
		subpath,
		&options)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(string(commitDir), 0755)
	if err != nil {
		return nil, err
	}

	for _, file := range folder {
		if string(file.GetType()) == "dir" {
			err = os.MkdirAll(string(commitDir+file.GetPath()), 0755)
			if err != nil {
				return nil, err
			} else {
				fmt.Println("Folder created:", file.GetPath())
			}
			_, err = downloadFolder(ownername, reponame, file.GetPath(), commitSha, branchRef, client)
			if err != nil {
				return nil, err
			}
		} else if string(file.GetType()) == "file" {
			err = downloadFile(file.GetDownloadURL(), commitDir+file.GetPath())
			if err != nil {
				return nil, err
			}
		}
	}
	return folder, err
}

//downloadFile downloads a single file.
func downloadFile(url string, filepath string) error {
	fmt.Println("downloadFile", url)
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
