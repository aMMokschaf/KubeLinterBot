//Package handleresult removes the files after KubeLinter is done linting. It passes a status back to main
//to decide if a comment will be posted or not.
package handleresult

import (
	"fmt"
	"main/internal/parsehook"
	"main/internal/postcomment"
	"os"
	"path/filepath"
)

//Handle calls removeDownloadedFiles after linting. After this, it passes kubelinters exit-code back.
func Handle(data parsehook.ParseResult, result []byte, status error, dir string) error {
	err := RemoveDownloadedFiles("./downloadedYaml/"+dir+"/", 0)
	fmt.Println("Removing downloaded files after linting...")
	if err != nil {
		fmt.Println("Error while removing files:\n", err)
	} else {
		fmt.Println("Files removed.")
	}
	if status != nil {
		if data.Event == "push" {
			err = postcomment.Push(data.Push.OwnerName, data.Push.RepoName, data.Push.Sha, result)
		} else if data.Event == "pull" {
			err = postcomment.PullRequestReview(data.Pull.OwnerName, data.Pull.RepoName, data.Pull.Sha, data.Pull.Number, result)
		}
		if err != nil {
			return err
		}
	} else {
		return nil
	}
	return nil
}

//RemoveDownloadedFiles removes all downloaded files in order to keep the storage-requirements low.
func RemoveDownloadedFiles(dir string, debug int) error {
	if debug == 1 {
		return nil
	}
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
