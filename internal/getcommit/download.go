//Package getcommit is used to download all folders with .yaml and .yml-files.
package getcommit

import (
	"context"
	"fmt"
	"io"
	"main/internal/authentication"
	"os"
	"regexp"

	"github.com/google/go-github/github"
)

const mainDir = "./downloadedYaml/"

//DownloadCommit downloads all .yaml or .yml-files.
//These are then passed to the KubeLinter-binary.
func DownloadCommit(ownername string, reponame string, commitSha string, branch string, filenames []string, number int, client authentication.Client) ([]*github.RepositoryContent, error) {
	err2 := download(ownername, reponame, "", commitSha, branch, filenames, number, client)
	if err2 != nil {
		return nil, err2
	}
	return nil, nil
}

func download(ownername string, reponame string, subpath string, commitSha string, branch string, filenames []string, number int, client authentication.Client) error {
	fmt.Println("Downloading contents...")
	//fmt.Println(ownername, reponame, subpath, commitSha, branch, filenames, number, client)
	var downloadDir string
	if commitSha != "" {
		downloadDir = mainDir + commitSha + "/"
	} else {
		downloadDir = mainDir + "_" + ownername + "_" + reponame
	}
	fmt.Println("Downloaddir:", downloadDir)

	branchRef, _, err := client.GithubClient.Repositories.GetBranch(context.Background(), ownername, reponame, branch)
	if err != nil {
		fmt.Println("getbranch", err)
		return err
	}

	var options = github.RepositoryContentGetOptions{Ref: branchRef.GetName()}
	for _, file := range filenames {
		f, err := client.GithubClient.Repositories.DownloadContents(context.Background(), ownername, reponame, file, &options)
		if err != nil {
			fmt.Println("Error downloadcontents", err)
			return nil
		} else {
			err := downloadSingleFile(f, downloadDir, file)
			if err != nil {
				fmt.Println("Error downloadSingleFile", err)
			}
		}
	}
	return nil
}

func downloadSingleFile(data io.ReadCloser, downloadDir string, filename string) error {
	fmt.Println("Downloading file:", filename)

	re := regexp.MustCompile(".*/")
	path := re.FindStringSubmatch(filename)

	if len(path) == 0 {
		err := os.MkdirAll(downloadDir, 0755)
		if err != nil {
			return err
		}
	} else {
		err := os.MkdirAll(downloadDir+path[0], 0755)
		if err != nil {
			return err
		}
	}

	out, err := os.Create(downloadDir + filename)
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
