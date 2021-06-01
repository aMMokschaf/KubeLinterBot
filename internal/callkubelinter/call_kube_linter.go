//Package callkubelinter checks if the KubeLinter-binary exists, updates and executes it.
package callkubelinter

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

//CallKubelinter calls the kube-linter binary in kubelinter/-folder.
func CallKubelinter(yamlDir string) ([]byte, error) {
	cmd := exec.Command("kubelinter/kube-linter", "lint", yamlDir)
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
	st, err := os.Stat("./kubelinter/kube-linter")

	if err != nil {
		fmt.Println("Could not find KubeLinter. Please download the latest release: https://github.com/stackrox/kube-linter/releases", err)

		return err
	}

	if st.Mode()&0111 == 0 {
		return errors.New("kube-linter binary found, but is not executable")
	}
	fmt.Println("KubeLinter found.")
	return nil
}

//export after implementation. Security issues? Maybe only check for new version, no download.
func checkForKubeLinterUpdate() {
	//get KubeLinter-Version
	//compare to latest release
	//update: automatically? How to make sure that linting is not affected?
}
