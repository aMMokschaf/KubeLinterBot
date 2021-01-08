package callkubelinter

import (
	"fmt"
	"os/exec"
)

//TODO: Doc
func Callkubelinter() []byte {
	//TODO: Change folder?
	cmd := exec.Command("kubelinter/kube-linter", "lint", "./downloadedYaml/")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Kube-Linter has found problems with your files.\n", err)
		fmt.Printf("%s\n", out)
	} else {
		fmt.Println("Kube-Linter has not found any problems with your files.")
	}
	return out
}
