//Package callkubelinter checks if the KubeLinter-binary exists, updates and executes it.
package callkubelinter

import (
	"fmt"
	"os"
	"os/exec"
)

//CallKubelinter calls the kube-linter binary in kubelinter/-folder.
func CallKubelinter() ([]byte, error) {
	checkForKubeLinterBinary()
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

//export und woanders aufrufen
func checkForKubeLinterBinary() {
	f, err := os.Open("./kubelinter/kube-linter")
	if err != nil {
		fmt.Println("Could not find KubeLinter.", err)
	} else {
		fmt.Println("KubeLinter found.")
	}
	defer f.Close()
}

//export after implementation
func checkForKubeLinterUpdate() {
	//get KubeLinter-Version
	//compare to latest release
	//update: automatically? How to make sure that linting is not affected?

}
