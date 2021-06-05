// Package handleresult removes the files after KubeLinter is done linting. It passes a status back to main
// to decide if a comment will be posted or not.
package handleresult

import (
	"fmt"
	"os"

	"github.com/aMMokschaf/KubeLinterBot/internal/authentication"
	"github.com/aMMokschaf/KubeLinterBot/internal/parsehook"
	"github.com/aMMokschaf/KubeLinterBot/internal/postcomment"

	"path/filepath"
)

// Handle calls removeDownloadedFiles after linting. After this, it passes kubelinters exit-code back.
func Handle(data *parsehook.GeneralizedResult, kubeLinterOutput []byte, status error, dir string, client *authentication.Client) error {
	err := RemoveDownloadedFiles(dir)
	if err != nil {
		fmt.Println("Error while removing files:\n", err)
	} else {
		fmt.Println("Files removed.")
	}
	if status == nil {
		return nil
	}
	err = postcomment.PostComment(data, kubeLinterOutput, client)
	return err
}

// RemoveDownloadedFiles removes all downloaded files in order to keep the storage-requirements low.
func RemoveDownloadedFiles(dir string) error {
	fmt.Println("Removing downloaded files after linting...")
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
