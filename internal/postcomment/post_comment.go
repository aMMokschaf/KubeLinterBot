//Package postcomment handles the posting of comments with KubeLinter's linting-results to the appropriate commit.
package postcomment

import (
	"context"
	"fmt"
	"main/internal/authentication"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

//PostCommentPush authorizes with github to post KubeLinter's results to the commit.
func PostCommentPush(username string, reponame string, commitSha string, result []byte) error {
	githubClient := authentication.GetGithubClient()

	var bdy string = string(result)
	comment := github.RepositoryComment{Body: &bdy}
	_, _, err := githubClient.Repositories.CreateComment(oauth2.NoContext, username, reponame, commitSha, &comment)
	if err != nil {
		fmt.Println("Posting Kubelinter's comment failed, error:", err)
		return err
	} else {
		fmt.Println("Comment posted successfully.\nKubeLinterBot is listening for Webhooks...")
		return nil
	}
}

//PostPullRequestReviewWithComment TODO blabla
func PostPullRequestReviewWithComment(username string, reponame string, commitSha string, result []byte) error {
	fmt.Println("postpullRequestReviewWithComment method")
	fmt.Println(username, reponame, commitSha, result)
	githubClient := authentication.GetGithubClient()
	//var PullRequestReviewRequest
	var comments []*github.DraftReviewComment
	var comment github.DraftReviewComment
	var commentString = "blablatest"
	fmt.Println(commentString)
	comment.Body = &commentString

	var review github.PullRequestReviewRequest

	review.Comments = comments

	rev, resp, err := githubClient.PullRequests.CreateReview(context.Background(), "aMMokschaf", "yamls", 36, &review)
	fmt.Println("pullrequestreview", rev, resp, err)
	rev, resp, err = githubClient.PullRequests.SubmitReview(context.Background(), "aMMokschaf", "yamls", 36, 0, &review)
	fmt.Println("pullrequestreview", rev, resp, err)

	//---

	// fmt.Println(username, reponame, commitSha, result)
	// githubClient := authentication.GetGithubClient()

	// var commentString = "blablatest"
	// fmt.Println(commentString)

	// var review github.PullRequestReview

	// review.Body = &commentString
	// rev, _, err := githubClient.PullRequests.CreateReview(context.Background(), "KubeLinterBot", "yamls", 17, review)
	// fmt.Println("pullrequestreview", rev, err)

	//---

	//githubClient.PullRequests.SubmitReview(context.Background(), "aMMokschaf", "yamls", 17, 0, rev)

	//githubClient.PullRequests.SubmitReview(context.Background(), ownername, reponame, pullnumber)

	// githubClient := authentication.GetGithubClient()

	// var bdy string = string(result)
	// // comment := github.RepositoryComment{Body: &bdy}
	// // _, _, err := githubClient.Repositories.CreateComment(oauth2.NoContext, username, reponame, commitSha, &comment)
	// comment = github.DraftReviewComment("", )
	// if err != nil {
	// 	fmt.Println("Posting Kubelinter's comment failed, error:", err)
	// 	return err
	// } else {
	// 	fmt.Println("Comment posted successfully.\nKubeLinterBot is listening for Webhooks...")
	// 	return nil
	// }
	return nil
}
