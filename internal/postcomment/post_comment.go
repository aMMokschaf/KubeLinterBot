//Package postcomment handles the posting of comments with KubeLinter's linting-results to the appropriate commit.
package postcomment

import (
	"context"
	"fmt"
	"main/internal/authentication"
	"strings"

	"github.com/google/go-github/github"
)

//Push authorizes with github to post KubeLinter's results to the commit.
func Push(username string, reponame string, commitSha string, result []byte, client *authentication.Client) error {
	cleanResult, numberOfLintErrors := separateComments(result, commitSha)
	fmt.Println("clean:", cleanResult, numberOfLintErrors)

	var bdy string = string(result)
	comment := github.RepositoryComment{Body: &bdy}
	_, _, err := client.GithubClient.Repositories.CreateComment(context.Background(), username, reponame, commitSha, &comment)
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
	var comments []*github.DraftReviewComment

	cleanResult, commentNumberOfLintErrors := separateComments(result, commitSha)

	for _, file := range files {
		var commentPath = file
		var commentPosition = 1 //indicating line 1
		var commentBody string
		for i := 0; i < len(cleanResult); i++ {
			if strings.Contains(cleanResult[i], file) {
				commentBody += cleanResult[i] + "\n\n"
			}
		}
		var comment = github.DraftReviewComment{}
		comment.Path = &commentPath
		comment.Position = &commentPosition
		comment.Body = &commentBody

		comments = append(comments, &comment)
	}

	body := "KubeLinter has found possible security- or production-readiness-errors. Please check the comments made by KubeLinterBot for each file.\n\n" + commentNumberOfLintErrors
	event := "REQUEST_CHANGES"
	var review = github.PullRequestReviewRequest{
		Body:     &body,
		Event:    &event,
		Comments: comments,
	}

	re, resp, err := client.GithubClient.PullRequests.CreateReview(context.Background(), username, reponame, int(number), &review)
	fmt.Println("create review re", re, "\nresp", resp, "\nerr", err)

	return nil
}

func separateComments(result []byte, commitSha string) ([]string, string) {
	comments := strings.Split(string(result), "\n")
	for i, comment := range comments {
		if comment != "" {
			comments[i] = cleanUpComment(comment, commitSha)
		}
	}
	fmt.Println("End of separate:", comments[len(comments)-2])
	return comments, comments[len(comments)-2]
}

func cleanUpComment(comment string, commitSha string) string {
	fmt.Println("comment vorher:", comment)
	comment = strings.Replace(comment, "downloadedYaml/"+commitSha+"/", "", 1)
	fmt.Println("comment nachher:", comment)
	return comment
}
