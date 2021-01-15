package callkubelinter

import (
	"fmt"
	"os/exec"
)

//Calls the kube-linter binary in kubelinter/-folder.
func Callkubelinter() ([]byte, int) {
	//TODO: Change folder?
	cmd := exec.Command("kubelinter/kube-linter", "lint", "./downloadedYaml/")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Kube-Linter has found problems with your files.\n", err)
		fmt.Printf("%s\n", out)
		return out, 1
	} else {
		fmt.Println("Kube-Linter has not found any problems with your files.")
		return out, 0
	}
}
