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

//DownloadCommit authenticates with oauth and downloads all folders with .yaml or .yml-files.
//These are then passed to the KubeLinter-binary.
func DownloadCommit(ownername string, reponame string, commitSha string, branch string, filenames []string, number int, client authentication.Client) ([]*github.RepositoryContent, error) {
	fmt.Println("DownloadCommit")
	fmt.Println(ownername, reponame, commitSha, branch, filenames, number, client)
	// repoContent, err := downloadFolder(ownername, reponame, "", commitSha, branch, client.GithubClient)
	// if err != nil {
	// 	return nil, err
	// }
	err2 := download(ownername, reponame, "", commitSha, branch, filenames, number, client)
	if err2 != nil {
		return nil, err2
	}
	// return repoContent, err
	return nil, nil
}

func download(ownername string, reponame string, subpath string, commitSha string, branch string, filenames []string, number int, client authentication.Client) error {
	fmt.Println("Downloading contents...")
	fmt.Println(ownername, reponame, subpath, commitSha, branch, filenames, number, client)
	var downloadDir string
	if commitSha != "" {
		downloadDir = mainDir + commitSha + "/"
	} else {
		downloadDir = mainDir + "_" + ownername + "_" + reponame
	}
	fmt.Println("Downloaddir:", downloadDir)

	// err := os.Mkdirall(downloadDir, 0755) //TODO check permissions
	// if err != nil {
	// 	return err
	// }

	//branchRef statt _
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
			//do something
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
	fmt.Println(path)
	if len(path) == 0 {
		err := os.MkdirAll(downloadDir, 0755)
		if err != nil {
			fmt.Println("Error while os.MkDirall()", err)
			return err
		}
	} else {
		err := os.MkdirAll(downloadDir+path[0], 0755)
		if err != nil {
			fmt.Println("Error while os.MkDirall()", err)
			return err
		}
	}

	out, err := os.Create(downloadDir + filename)
	if err != nil {
		fmt.Println("Error while os.Create()", err)
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, data)
	if err != nil {
		fmt.Println("Error while io.Copy()", err)
		return err
	}
	return nil
}

// //downloadFolder downloads all files in a folder, creating subfolders as necessary.
// func downloadFolder(ownername string, reponame string, subpath string, commitSha string, branchRef string, client *github.Client) ([]*github.RepositoryContent, error) {
// 	fmt.Println("downloadFolder")
// 	fmt.Println(ownername, reponame, subpath, commitSha, branchRef, client)
// 	var commitDir string
// 	if commitSha != "" {
// 		commitDir = mainDir + commitSha + "/"
// 	} else {
// 		commitDir = mainDir + ownername + reponame + "/"
// 	}

// 	branch, _, err := client.Repositories.GetBranch(context.Background(),
// 		ownername,
// 		reponame,
// 		branchRef)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var options = github.RepositoryContentGetOptions{Ref: branch.GetName()}
// 	//f, err := client.Repositories.DownloadContents(context.Background(), ownername, reponame, subpath, &options)

// 	_, folder, _, err := client.Repositories.GetContents(context.Background(),
// 		ownername,
// 		reponame,
// 		subpath,
// 		&options)
// 	if err != nil {
// 		return nil, err
// 	}

// 	err = os.MkdirAll(string(commitDir), 0755)
// 	if err != nil {
// 		return nil, err
// 	}

// 	for _, file := range folder {
// 		if string(file.GetType()) == "dir" {
// 			err = os.MkdirAll(string(commitDir+file.GetPath()), 0755)
// 			if err != nil {
// 				return nil, err
// 			} else {
// 				fmt.Println("Folder created:", file.GetPath())
// 			}
// 			_, err = downloadFolder(ownername, reponame, file.GetPath(), commitSha, branchRef, client)
// 			if err != nil {
// 				return nil, err
// 			}
// 		} else if string(file.GetType()) == "file" {
// 			err = downloadFile(file.GetDownloadURL(), commitDir+file.GetPath())
// 			if err != nil {
// 				return nil, err
// 			}
// 		}
// 	}
// 	return folder, err
// }

// //downloadFile downloads a single file.
// func downloadFile(url string, filepath string) error {
// 	fmt.Println("downloadFile", url)
// 	out, err := os.Create(filepath)
// 	if err != nil {
// 		return err
// 	}
// 	defer out.Close()

// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	_, err = io.Copy(out, resp.Body)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
