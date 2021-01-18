//Package postcomment handles the posting of comments with KubeLinter's linting-results to the appropriate commit.
package postcomment

import (
	"fmt"
	"main/internal/authentication"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

//PostComment authorizes with github to post KubeLinter's results to the commit.
func PostComment(username string, reponame string, commitSha string, result []byte) error {
	githubClient := authentication.GetGithubClient()

	var bdy string = string(result)
	comment := github.RepositoryComment{Body: &bdy}
	_, _, err := githubClient.Repositories.CreateComment(oauth2.NoContext, username, reponame, commitSha, &comment)
	if err != nil {
		fmt.Println("Posting Kubelinter's comment failed, error:", err)
		return err
	} else {
		fmt.Println("Comment posted successfully.\nKubeLinterBot is listening for Webhooks.")
		return nil
	}
}
