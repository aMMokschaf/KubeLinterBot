package getcommit

import (
	"fmt"
)

func GetCommit(token string, username string, reponame string, commitSha string, addedFiles []string, modifiedFiles []string) {
	fmt.Println("get commit method")
	fmt.Println(addedFiles, "\n", modifiedFiles)

	//TODO remove hardcoded stuff
	DownloadCommit(token, username, reponame, commitSha, addedFiles, modifiedFiles)
}
