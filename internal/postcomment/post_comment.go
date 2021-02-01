//Package postcomment handles the posting of comments with KubeLinter's linting-results to the appropriate commit.
package postcomment

import (
	"context"
	"fmt"
	"main/internal/authentication"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

//Push authorizes with github to post KubeLinter's results to the commit.
func Push(username string, reponame string, commitSha string, result []byte, client *authentication.Client) error {
	//githubClient := authentication.GetGithubClient()

	var bdy string = string(result)
	comment := github.RepositoryComment{Body: &bdy}
	_, _, err := client.GithubClient.Repositories.CreateComment(oauth2.NoContext, username, reponame, commitSha, &comment)
	if err != nil {
		fmt.Println("Posting Kubelinter's comment failed, error:", err)
		return err
	} else {
		fmt.Println("Comment posted successfully.\nKubeLinterBot is listening for Webhooks...")
		return nil
	}
}

//PullRequestReview TODO blabla
func PullRequestReview(username string, reponame string, commitSha string, number int64, result []byte, client *authentication.Client) error {
	fmt.Println("postpullRequestReviewWithComment method")
	fmt.Println(username, reponame, commitSha, string(result))
	//githubClient := authentication.GetGithubClient()

	//Kommentar erstellen mit REQUEST CHANGES DraftReviewComment

	var commentPath = "pod.yaml"
	var commentPosition = 1 //indicating line 1
	var commentBody = "blablatest"
	fmt.Println(commentPath, commentPosition, commentBody)

	var comment = github.DraftReviewComment{}

	comment.Path = &commentPath
	comment.Position = &commentPosition
	comment.Body = &commentBody
	fmt.Println("comment", comment)

	var comments []*github.DraftReviewComment
	comments = append(comments, &comment)
	fmt.Println(comments)

	var stringEvent = "REQUEST_CHANGES"
	// var commitIDString = "" //empty weil dann neuester commit //commitSha

	//var review2 = github.PullRequestReviewRequest{CommitID: &commitIDString, Body: &commentBody, Event: &stringEvent, Comments: comments}
	// fmt.Println(review)
	body := "Ja das ist hier!"
	event := "COMMENT"
	var review = github.PullRequestReviewRequest{
		Body:  &body,
		Event: &event,
	} //{Body: &commentBody, Event: &stringEvent}

	// u := fmt.Sprintf("repos/%v/%v/pulls/%d/comments", "aMMokschaf", "yamls", 39)

	// req, err := githubClient.NewRequest("POST", u, review)
	// if err != nil {
	// 	fmt.Println("newRequest", err)
	// 	return err
	// }

	// r := new(github.PullRequestReview)
	// resp, err := githubClient.Do(context.Background(), req, r)
	// if err != nil {
	// 	fmt.Println("Do", err)
	// 	return err
	//}

	re, resp, err := client.GithubClient.PullRequests.CreateReview(context.Background(), username, reponame, int(number), &review)
	fmt.Println("create review re", re, "\nresp", resp, "\nerr", err)

	review2 := github.PullRequestReviewRequest{}

	review2.Body = &commentBody
	review2.Event = &stringEvent
	re, resp, err = client.GithubClient.PullRequests.SubmitReview(context.Background(), username, reponame, int(number), re.GetID(), &review2)
	fmt.Println("submit review re", re, "\nresp", resp, "\nerr", err)
	return nil
}

// //PullRequestReview TODO blabla
// func PullRequestReview(username string, reponame string, commitSha string, result []byte) error {
// 	fmt.Println("postpullRequestReviewWithComment method")
// 	fmt.Println(username, reponame, commitSha, string(result))
// 	githubClient := authentication.GetGithubClient()
// 	//var PullRequestReviewRequest

// 	var commentPath = "server-deployment/deeper/deep-server.yaml"
// 	var commentPosition = 0
// 	var commentBody = "blablatest"

// 	var comment = github.DraftReviewComment{}

// 	fmt.Println(commentPath, commentPosition, commentBody)
// 	comment.Path = &commentPath
// 	//comment.Position = &commentPosition
// 	comment.Body = &commentBody
// 	fmt.Println(comment)

// 	var comments []*github.DraftReviewComment
// 	comments = append(comments, &comment)
// 	fmt.Println(comments)

// 	var stringEvent = "COMMENT"
// 	var commitIDString = ""

// 	//var review = github.PullRequestReviewRequest{CommitID: &commitIDString, Body: &commentBody, Event: &stringEvent, Comments: comments}
// 	var review = github.PullRequestReviewRequest{CommitID: &commitIDString, Body: &commentBody, Event: &stringEvent, Comments: nil}
// 	fmt.Println(review)
// 	// review.Event = &stringEvent
// 	// review.CommitID = &stringEvent

// 	// review.Comments = comments

// 	// var id = int64(77734728)
// 	// var userLogin = "KubeLinterBot"
// 	// var user = github.User{Login: &userLogin, ID: &id}
// 	// var users []*github.User
// 	// append(users, &user)
// 	// var reviewer = github.Reviewers{Users: users}
// 	var reviewers []string
// 	reviewers = append(reviewers, "KubeLinterBot")
// 	// var revReq = github.ReviewersRequest{Reviewers: reviewers}
// 	// pr, r, err := githubClient.PullRequests.RequestReviewers(context.Background(), "aMMokschaf", "yamls", 37, revReq)
// 	// fmt.Println("pr", pr, "r", r, "err", err)
// 	rev, resp, err := githubClient.PullRequests.CreateReview(context.Background(), "aMMokschaf", "yamls", 39, &review)
// 	fmt.Println("create review rev", rev, "\nresp", resp, "\nerr", err)

// 	// var grev = github.PullRequestReview{}
// 	// grev.Body = &commentBody
// 	// grev

// 	// githubClient.PullRequests.SubmitReview(context.Background(), "aMMokschaf", "yamls", 17, 0, rev)

// 	// rev, resp, err = githubClient.PullRequests.SubmitReview(context.Background(), "aMMokschaf", "yamls", 36, 0, &review)
// 	// fmt.Println("submit review rev", rev, "\nresp", resp, "\nerr", err)

// 	//---

// 	// fmt.Println(username, reponame, commitSha, result)
// 	// githubClient := authentication.GetGithubClient()

// 	// var commentString = "blablatest"
// 	// fmt.Println(commentString)

// 	// var review github.PullRequestReview

// 	// review.Body = &commentString
// 	// rev, _, err := githubClient.PullRequests.CreateReview(context.Background(), "KubeLinterBot", "yamls", 17, review)
// 	// fmt.Println("pullrequestreview", rev, err)

// 	//---

// 	//githubClient.PullRequests.SubmitReview(context.Background(), "aMMokschaf", "yamls", 17, 0, rev)

// 	//githubClient.PullRequests.SubmitReview(context.Background(), ownername, reponame, pullnumber)

// 	// githubClient := authentication.GetGithubClient()

// 	// var bdy string = string(result)
// 	// // comment := github.RepositoryComment{Body: &bdy}
// 	// // _, _, err := githubClient.Repositories.CreateComment(oauth2.NoContext, username, reponame, commitSha, &comment)
// 	// comment = github.DraftReviewComment("", )
// 	// if err != nil {
// 	// 	fmt.Println("Posting Kubelinter's comment failed, error:", err)
// 	// 	return err
// 	// } else {
// 	// 	fmt.Println("Comment posted successfully.\nKubeLinterBot is listening for Webhooks...")
// 	// 	return nil
// 	// }
// 	return nil
// }
