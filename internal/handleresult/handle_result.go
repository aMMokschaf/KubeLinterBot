package handleresult

import (
	"fmt"
	"os"
	"path/filepath"
)

//TODO: Doc, ausprogrammieren
func HandleResult() {
	fmt.Println("Entering HandleResult")
	//exit code 1: comment
	//exit code 0: no comment
	err := removeDownloadedFiles("./downloadedYaml/")
	if err != nil {
		fmt.Println("Error while removing files.")
	} else {
		fmt.Println("Files removed.")
	}
}

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
