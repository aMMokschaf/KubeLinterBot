package getcommit

import (
	"github.com/aMMokschaf/KubeLinterBot/internal/authentication"
	"github.com/aMMokschaf/KubeLinterBot/internal/parsehook"
)

//GetCommit TODO Kommentar
func GetCommit(result *parsehook.GeneralizedResult, client authentication.Client) (string, error) {
	return DownloadCommit(result.OwnerName,
		result.RepoName,
		result.Sha,
		result.Branch,
		result.AddedOrModifiedFiles,
		result.Number,
		client)
}
