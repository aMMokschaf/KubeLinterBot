package getcommit

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadFile(url string, filename string) error {
	//TODO implement subfolders
	fmt.Println("Downloading file " + filename + "\n")
	const folder = "./downloadedYaml/"
	out, err := os.Create(folder + filename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url + filename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
