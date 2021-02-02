package getcommit

import (
	"main/internal/authentication"
	"main/internal/parsehook"
)

//GetCommit TODO Kommentar
func GetCommit(result *parsehook.GeneralizedResult, client authentication.Client) {
	DownloadCommit(result.OwnerName,
		result.RepoName,
		result.Sha,
		result.Branch,
		result.AddedOrModifiedFiles,
		result.Number,
		client)
}
