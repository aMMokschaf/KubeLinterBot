package getcommit

import "main/internal/parsehook"

//GetCommit TODO Kommentar
func GetCommit(result parsehook.ParseResult) {
	var ownername string
	var reponame string
	var commitSha string
	var branchRef string
	/* to be used later
	var addedFiles []string
	var modifiedFiles []string
	*/

	if result.Event == "push" {
		//Push
		ownername = result.Push.UserName
		reponame = result.Push.RepoName
		commitSha = result.Push.Sha
		branchRef = result.Push.Branch
		/*to be used later
		addedFiles = result.Push.AddedFiles
		modifiedFiles = result.Push.ModifiedFiles
		*/

	} else if result.Event == "pull" {
		//Pull
		ownername = result.Pull.OwnerName
		reponame = result.Pull.RepoName
		commitSha = result.Pull.Sha
		branchRef = result.Pull.Branch
	}

	DownloadCommit(ownername, reponame, commitSha, branchRef)
}
