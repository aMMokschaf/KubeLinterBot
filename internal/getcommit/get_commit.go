package getcommit

import (
	"fmt"
)

func GetCommit(token string, username string, reponame string, commitSha string, addedFiles []string, modifiedFiles []string) {
	fmt.Println("Getting commit", commitSha, "from repository", reponame)
	fmt.Printf("Added Files:%v\nModified Files:%v\n", addedFiles, modifiedFiles)

	DownloadCommit(token, username, reponame, commitSha, addedFiles, modifiedFiles)
}
