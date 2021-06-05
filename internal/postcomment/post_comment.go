// Package postcomment handles the posting of comments with KubeLinter's linting-results to the appropriate commit.
package postcomment

import (
	"context"
	"fmt"
	"strings"

	"github.com/aMMokschaf/KubeLinterBot/internal/authentication"
	"github.com/aMMokschaf/KubeLinterBot/internal/parsehook"

	"github.com/google/go-github/github"
)

// PostComment pre-processes the output of KubeLinter and chooses the correct comment-method for commits and pull-requests.
func PostComment(data *parsehook.GeneralizedResult, kubelinterOutput []byte, client *authentication.Client) error {
	separatedComments, errorCountString := separateComments(kubelinterOutput, data.Sha)

	var err error
	if data.IsPush() {
		err = push(data.OwnerName, data.RepoName, data.Sha, kubelinterOutput, client, separatedComments, errorCountString)
	} else {
		err = pullRequestReview(data.BaseOwnerName, data.BaseRepoName, data.Sha, data.Number, data.AddedOrModifiedFiles, kubelinterOutput, client, separatedComments, errorCountString)
	}
	return err
}

// push is used to post a comment to a whole commit.
func push(username string, reponame string, commitSha string, result []byte, client *authentication.Client, cleanResult []string, numberOfLintErrors string) error {
	var bdy string = strings.Join(cleanResult, "\n")
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

// pullRequestReview is used to post a review for a Pull-Request.
func pullRequestReview(username string, reponame string, commitSha string, number int, files []string, result []byte, client *authentication.Client, separatedComments []string, errorCountString string) error {
	var comments []*github.DraftReviewComment

	for _, file := range files {
		var commentPath = file
		var commentPosition = 1 //indicating line 1, to be changed in future releases to accomodate real line numbers
		var commentBody string
		for i := 0; i < len(separatedComments); i++ {
			if strings.Contains(separatedComments[i], file) {
				commentBody += separatedComments[i] + "\n\n"
			}
		}
		var comment = github.DraftReviewComment{}
		comment.Path = &commentPath
		comment.Position = &commentPosition
		comment.Body = &commentBody

		comments = append(comments, &comment)
	}

	body := "KubeLinter has found possible security- or production-readiness-errors. Please check the comments made by KubeLinterBot for each file.\n\n" + errorCountString
	event := "REQUEST_CHANGES"
	var review = github.PullRequestReviewRequest{
		Body:     &body,
		Event:    &event,
		Comments: comments,
	}

	_, _, err := client.GithubClient.PullRequests.CreateReview(context.Background(), username, reponame, int(number), &review)
	if err != nil {
		return err
	}

	return nil
}

// separateComments splits Kubelinter's output to separate error messages.
func separateComments(result []byte, commitSha string) ([]string, string) {
	comments := strings.Split(string(result), "\n")
	for i, comment := range comments {
		if comment != "" {
			comments[i] = cleanUpComment(comment, commitSha)
		}
	}
	return comments, comments[len(comments)-2]
}

// cleanUpComment removes the leading "downloadedYaml/[commitSha]/ that is not part of the linted repository."
func cleanUpComment(comment string, commitSha string) string {
	comment = strings.Replace(comment, "downloadedYaml/"+commitSha+"/", "", 1)
	return comment
}
