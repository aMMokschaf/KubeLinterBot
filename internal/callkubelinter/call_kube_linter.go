package callkubelinter

import (
	"fmt"
	"os/exec"
)

func Callkubelinter() {
	//TODO take arguments
	output, err := exec.Command("kubelinter/kube-linter", "lint", "pod.yaml").Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(output))
}
