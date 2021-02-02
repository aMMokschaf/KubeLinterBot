//Package callkubelinter checks if the KubeLinter-binary exists, updates and executes it.
package callkubelinter

import (
	"fmt"
	"os"
	"os/exec"
)

//CallKubelinter calls the kube-linter binary in kubelinter/-folder.
func CallKubelinter() ([]byte, error) {
	//TODO: Change folder?
	cmd := exec.Command("kubelinter/kube-linter", "lint", "./downloadedYaml/")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Kube-Linter has found problems with your files:", err)
		fmt.Printf("%s\n", out)
	} else {
		fmt.Println("Kube-Linter has not found any problems with your files.")
	}
	return out, err
}

//CheckForKubeLinterBinary checks if a Kubelinter-binary exists in /Kubelinterbot/kubelinter/.
func CheckForKubeLinterBinary() error {
	f, err := os.Open("./kubelinter/kube-linter")
	defer f.Close()
	if err != nil {
		fmt.Println("Could not find KubeLinter. Please download the latest release: https://github.com/stackrox/kube-linter/releases", err)
		return err
	} else {
		fmt.Println("KubeLinter found.")
		return nil
	}
}

//export after implementation. Security issues? Maybe only check for new version, no download.
func checkForKubeLinterUpdate() {
	//get KubeLinter-Version
	//compare to latest release
	//update: automatically? How to make sure that linting is not affected?
}
