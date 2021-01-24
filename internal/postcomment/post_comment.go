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

	var commentPath = "server-deployment/deeper/deep-server.yaml"
	var commentPosition = 0
	var commentBody = "blablatest"

	var comment = github.DraftReviewComment{}

	fmt.Println(commentPath, commentPosition, commentBody)
	comment.Path = &commentPath
	//comment.Position = &commentPosition
	comment.Body = &commentBody
	fmt.Println(comment)

	var comments []*github.DraftReviewComment
	comments = append(comments, &comment)
	fmt.Println(comments)

	var stringEvent = "COMMENT"
	var commitIDString = ""

	//var review = github.PullRequestReviewRequest{CommitID: &commitIDString, Body: &commentBody, Event: &stringEvent, Comments: comments}
	var review = github.PullRequestReviewRequest{CommitID: &commitIDString, Body: &commentBody, Event: &stringEvent, Comments: nil}
	fmt.Println(review)
	// review.Event = &stringEvent
	// review.CommitID = &stringEvent

	// review.Comments = comments

	// var id = int64(77734728)
	// var userLogin = "KubeLinterBot"
	// var user = github.User{Login: &userLogin, ID: &id}
	// var users []*github.User
	// append(users, &user)
	// var reviewer = github.Reviewers{Users: users}
	var reviewers []string
	reviewers = append(reviewers, "KubeLinterBot")
	// var revReq = github.ReviewersRequest{Reviewers: reviewers}
	// pr, r, err := githubClient.PullRequests.RequestReviewers(context.Background(), "aMMokschaf", "yamls", 37, revReq)
	// fmt.Println("pr", pr, "r", r, "err", err)
	rev, resp, err := githubClient.PullRequests.CreateReview(context.Background(), "aMMokschaf", "yamls", 37, &review)
	fmt.Println("create review rev", rev, "\nresp", resp, "\nerr", err)

	// var grev = github.PullRequestReview{}
	// grev.Body = &commentBody
	// grev

	// githubClient.PullRequests.SubmitReview(context.Background(), "aMMokschaf", "yamls", 17, 0, rev)

	// rev, resp, err = githubClient.PullRequests.SubmitReview(context.Background(), "aMMokschaf", "yamls", 36, 0, &review)
	// fmt.Println("submit review rev", rev, "\nresp", resp, "\nerr", err)

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
