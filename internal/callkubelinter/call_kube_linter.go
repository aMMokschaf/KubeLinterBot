//Package callkubelinter checks if the KubeLinter-binary exists, updates and executes it.
package callkubelinter

import (
	"fmt"
	"os"
	"os/exec"
)

//Callkubelinter calls the kube-linter binary in kubelinter/-folder.
func Callkubelinter() ([]byte, error) {
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

func checkForKubeLinterBinary() {
	f, err := os.Open("./kubelinter/kube-linter")
	if err != nil {
		fmt.Println("Could not find Kubelinter.", err)
	} else {
		fmt.Println("KubeLinter found.")
	}
	defer f.Close()
}

func checkForKubeLinterUpdate() {}
