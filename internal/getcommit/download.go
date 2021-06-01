//Package getcommit is used to download all folders with .yaml and .yml-files.
package getcommit

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/aMMokschaf/KubeLinterBot/internal/authentication"

	"github.com/google/go-github/github"
)

const mainDir = "./downloadedYaml/"

//DownloadCommit downloads all .yaml or .yml-files.
//These are then passed to the KubeLinter-binary.
func DownloadCommit(ownername string, reponame string, commitSha string, branch string, filenames []string, number int, client authentication.Client) (string, error) {
	dir, err := download(ownername, reponame, "", commitSha, branch, filenames, number, client)
	if err != nil {
		return "", err
	}
	return dir, nil
}

func download(ownername string, reponame string, subpath string, commitSha string, branch string, filenames []string, number int, client authentication.Client) (string, error) { // better: return downloadDir from this function and pass to caller
	fmt.Println("Downloading contents...")
	downloadDir := mainDir + commitSha + "/"
	fmt.Println("Downloaddir:", downloadDir)

	branchRef, _, err := client.GithubClient.Repositories.GetBranch(context.Background(), ownername, reponame, branch) // add context parameter
	if err != nil {
		fmt.Println("getbranch", err)
		return "", err
	}

	var options = github.RepositoryContentGetOptions{Ref: branchRef.GetName()}
	for _, file := range filenames {
		f, err := client.GithubClient.Repositories.DownloadContents(context.Background(), ownername, reponame, file, &options) // add context param
		if err != nil {
			return "", err
		} else {
			err := downloadSingleFile(f, downloadDir, file)
			if err != nil {
				return "", err
			}
		}
	}
	return downloadDir, nil
}

func downloadSingleFile(data io.ReadCloser, downloadDir string, filename string) error {
	fmt.Println("Downloading file:", filename)
	dir := filepath.FromSlash(path.Dir(filename))

	if dir == "" {
		err := os.MkdirAll(downloadDir, 0700)
		if err != nil {
			return err
		}
	} else {
		err := os.MkdirAll(filepath.Join(downloadDir, dir), 0700)
		if err != nil {
			return err
		}
	}

	out, err := os.Create(filepath.Join(downloadDir, filename))
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, data)
	if err != nil {
		return err
	}
	return nil
}
