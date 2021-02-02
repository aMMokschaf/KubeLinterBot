package getcommit

import (
	"fmt"
	"main/internal/authentication"
	"main/internal/parsehook"
)

//GetCommit TODO Kommentar
func GetCommit(result *parsehook.GeneralizedResult, client authentication.Client) {
	fmt.Println("GetCommit")
	fmt.Println(result)
	fmt.Println(client)

	DownloadCommit(result.OwnerName,
		result.RepoName,
		result.Sha,
		result.Branch,
		result.AddedOrModifiedFiles,
		result.Number,
		client)
}
