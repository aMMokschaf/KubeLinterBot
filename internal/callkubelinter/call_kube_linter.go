//Package callkubelinter executes the KubeLinter-binary.
package callkubelinter

import (
	"fmt"
	"os/exec"
)

//Callkubelinter calls the kube-linter binary in kubelinter/-folder.
func Callkubelinter() ([]byte, error) {
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
