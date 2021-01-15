package handleresult

import (
	"fmt"
	"os"
	"path/filepath"
)

//HandleResult calls removeDownloadedFiles after linting. After this, it passes kubelinters exit-code back.
func HandleResult(status int) int {
	fmt.Println("Entering HandleResult")
	err := removeDownloadedFiles("./downloadedYaml/")
	if err != nil {
		fmt.Println("Error while removing files.")
	} else {
		fmt.Println("Files removed.")
	}
	if status == 1 {
		return 1
	} else {
		return 0
	}
}

//removeDownloadedFiles removes all downloaded files in order to keep the storage-requirements low.
func removeDownloadedFiles(dir string) error {
	fmt.Println("Entering removeDownloadedFiles()")
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
