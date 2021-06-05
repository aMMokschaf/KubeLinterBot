// Package getcommit is used to download all folders with .yaml and .yml-files.
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

// DownloadCommit downloads all .yaml or .yml-files.
func DownloadCommit(ownername string, reponame string, commitSha string, branch string, filenames []string, number int, client authentication.Client) (string, error) {
	fmt.Println("Downloading contents...")
	downloadDir := mainDir + commitSha + "/"
	fmt.Println("Downloaddir:", downloadDir)

	branchRef, _, err := client.GithubClient.Repositories.GetBranch(context.Background(), ownername, reponame, branch) // add context param
	if err != nil {
		return "", err
	}

	var options = github.RepositoryContentGetOptions{Ref: branchRef.GetName()}
	for _, file := range filenames {
		fmt.Println("Downloading file:", file)
		f, err := client.GithubClient.Repositories.DownloadContents(context.Background(), ownername, reponame, file, &options) // add context param
		if err != nil {
			return "", err
		} else {
			err := writeFileToDisk(f, downloadDir, file)
			if err != nil {
				return "", err
			}
		}
	}
	return downloadDir, nil
}

func writeFileToDisk(data io.ReadCloser, downloadDir string, filename string) error {
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
