package getcommit

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

/*
func checkArgsAndDownload() {
	if len(os.Args) != 3 {
		fmt.Println("usage: download url filename")
		os.Exit(1)
	}
	url := os.Args[1]
	filename := os.Args[2]

	err := DownloadFile(url, filename)
	if err != nil {
		panic(err)
	}

}*/

// DownloadFile will download a url and store it in local filepath.
// It writes to the destination file as it downloads it, without
// loading the entire file into memory.
func DownloadFile(url string, filename string) error {
	// Create the file
	fmt.Println("Downloading file " + filename + "\n")
	const folder = "/downloadedYaml/"
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
