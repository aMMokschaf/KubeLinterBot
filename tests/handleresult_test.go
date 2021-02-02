//Package tests consists of multiple tests for critical components
package tests

import (
	"fmt"
	"io"
	"os"
	"testing"

	"main/internal/handleresult"
)

/*TestRemoveDownloadedFiles creates:

./downloadedYAML/testfile.go,
./downloadedYAML/firstfolder/testfile2.go
./downloadedYAML/firstfolder/secondfolder/testfile3.go,

and  deletes them afterwards. It then checks if ./downloadedYAML/ is empty.

*/
func TestRemoveDownloadedFiles(t *testing.T) {
	err := os.MkdirAll("../downloadedYaml/firstfolder/secondfolder", 0755)
	if err != nil {
		t.Error("Setting up TestRemoveDownloadedFiles failed.", err)
	}
	_, err = os.Create("../downloadedYaml/testfile.go")
	if err != nil {
		t.Error("Setting up TestRemoveDownloadedFiles failed.", err)
	}
	_, err = os.Create("../downloadedYaml/firstfolder/testfile2.go")
	if err != nil {
		t.Error("Setting up TestRemoveDownloadedFiles failed.", err)
	}
	_, err = os.Create("../downloadedYaml/firstfolder/secondfolder/testfile3.go")
	if err != nil {
		t.Error("Setting up TestRemoveDownloadedFiles failed.", err)
	}

	handleresult.RemoveDownloadedFiles("../downloadedYaml/", 0)
	f, err := os.Open("../downloadedYaml/")
	if err != nil {
		fmt.Println("Could not open folder.", err)
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err != io.EOF {
		t.Error("Directory is not empty.")
	}
}
