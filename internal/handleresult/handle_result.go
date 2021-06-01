//Package handleresult removes the files after KubeLinter is done linting. It passes a status back to main
//to decide if a comment will be posted or not.
package handleresult

import (
	"fmt"
	"os"

	"github.com/aMMokschaf/KubeLinterBot/internal/authentication"
	"github.com/aMMokschaf/KubeLinterBot/internal/parsehook"
	"github.com/aMMokschaf/KubeLinterBot/internal/postcomment"
)

//Handle calls removeDownloadedFiles after linting. After this, it passes kubelinters exit-code back.
func Handle(data *parsehook.GeneralizedResult, result []byte, status error, dir string, client *authentication.Client) error {
	err := RemoveDownloadedFiles(dir)
	if err != nil {
		fmt.Println("Error while removing files:\n", err)
	} else {
		fmt.Println("Files removed.")
	}
	if status == nil {
		return nil
	}
	// if data.Number == 0 { //isPush-methode bei GeneralizedResult einführen
	// 	err = postcomment.Push(data.OwnerName, data.RepoName, data.Sha, result, client)
	// } else {
	// 	err = postcomment.PullRequestReview(data.BaseOwnerName, data.BaseRepoName, data.Sha, data.Number, data.AddedOrModifiedFiles, result, client)
	// }
	err = postcomment.PostComment(data, result, client) //result in kubelinteroutput umbenennen
	return err
}

//RemoveDownloadedFiles removes all downloaded files in order to keep the storage-requirements low.
func RemoveDownloadedFiles(dir string) error {
	fmt.Println("Removing downloaded files after linting...")
	return os.RemoveAll(dir)
}
