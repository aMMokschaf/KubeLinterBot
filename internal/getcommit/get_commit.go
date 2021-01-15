package getcommit

import (
	"fmt"
)

func GetCommit(reponame string, addedFiles []string, modifiedFiles []string, token string) {
	fmt.Println("get commit method")
	fmt.Println(addedFiles, "\n", modifiedFiles)

	/*output, err := exec.Command("kubelinter/kube-linter", "lint", "pod.yaml").Output()
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(output))*/
	// for _, filename := range modifiedFiles {
	// 	DownloadFilePublic(reponame, filename)
	// }
	//TODO remove hardcoded stuff
	DownloadCommit(token, "aMMokschaf", "yamls", "f8a065ff6efd41d3212dee71f825d8094a155795", addedFiles, modifiedFiles)
}
