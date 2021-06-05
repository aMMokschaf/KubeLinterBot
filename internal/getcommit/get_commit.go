package getcommit

import (
	"github.com/aMMokschaf/KubeLinterBot/internal/authentication"
	"github.com/aMMokschaf/KubeLinterBot/internal/parsehook"
)

//GetCommit passes on the necessary data to download all needed files.
func GetCommit(result *parsehook.GeneralizedResult, client *authentication.Client) (string, error) {
	return DownloadCommit(result.OwnerName,
		result.RepoName,
		result.Sha,
		result.Branch,
		result.AddedOrModifiedFiles,
		result.Number,
		*client)
}
