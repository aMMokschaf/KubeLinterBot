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
func PullRequestReview(username string, reponame string, commitSha string, number int, files []string, result []byte, client *authentication.Client) error {
	fmt.Println("postpullRequestReviewWithComment method")
	fmt.Println(username, reponame, commitSha, number, files, string(result), client)

	var comments []*github.DraftReviewComment
	//build comments
	for _, file := range files {
		var commentPath = file
		var commentPosition = 1 //indicating line 1
		var commentBody = string(result)
		fmt.Println("Building comment:\n", commentPath, commentBody, commentPosition)
		var comment = github.DraftReviewComment{}
		comment.Path = &commentPath
		comment.Position = &commentPosition
		comment.Body = &commentBody

		//append comment to array of comments
		comments = append(comments, &comment)
	}

	fmt.Println(comments)

	//review erstellen

	body := "KubeLinter has found possible security- or production-readiness-errors. Please check the comments made by KubeLinterBot for each file."
	event := "COMMENT"
	var review = github.PullRequestReviewRequest{
		Body:     &body,
		Event:    &event,
		Comments: comments,
	}

	//review abschicken createReview()

	re, resp, err := client.GithubClient.PullRequests.CreateReview(context.Background(), username, reponame, int(number), &review)
	fmt.Println("create review re", re, "\nresp", resp, "\nerr", err)

	return nil
}
