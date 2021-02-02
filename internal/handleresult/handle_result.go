//Package handleresult removes the files after KubeLinter is done linting. It passes a status back to main
//to decide if a comment will be posted or not.
package handleresult

import (
	"fmt"
	"main/internal/authentication"
	"main/internal/parsehook"
	"main/internal/postcomment"
	"os"
	"path/filepath"
)

//Handle calls removeDownloadedFiles after linting. After this, it passes kubelinters exit-code back.
func Handle(data *parsehook.GeneralizedResult, result []byte, status error, dir string, client *authentication.Client) error {
	err := RemoveDownloadedFiles("./downloadedYaml/"+dir+"/", 0)
	if err != nil {
		fmt.Println("Error while removing files:\n", err)
	} else {
		fmt.Println("Files removed.")
	}
	if status != nil {
		if data.Number == 0 {
			err = postcomment.Push(data.OwnerName, data.RepoName, data.Sha, result, client)
		} else {
			err = postcomment.PullRequestReview(data.OwnerName, data.RepoName, data.Sha, data.Number, result, client)
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
	fmt.Println("Removing downloaded files after linting...")
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
	err = os.Remove(dir)
	if err != nil {
		return err
	}
	return nil
}
