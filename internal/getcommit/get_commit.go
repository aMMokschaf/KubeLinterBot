package getcommit

//GetCommit TODO Kommentar
func GetCommit(token string, ownername string, reponame string, commitSha string, branchRef string, addedFiles []string, modifiedFiles []string) {
	// if addedFiles != nil {
	// 	fmt.Println("Getting commit", commitSha, "from repository", reponame)
	// 	fmt.Printf("Added Files:%v\nModified Files:%v\n", addedFiles, modifiedFiles)
	// }

	DownloadCommit(ownername, reponame, commitSha, branchRef, addedFiles, modifiedFiles)
}
