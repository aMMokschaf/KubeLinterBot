// Package parsehook parses a github-webhook. Push-Events and Pull-Request-Events are handled.
package parsehook

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/aMMokschaf/KubeLinterBot/internal/authentication"

	"github.com/google/go-github/github"
	githubWebhook "gopkg.in/go-playground/webhooks.v5/github"
)

// GeneralizedResult is a representation of the fields of push- and pull-payloads that are
// interesting for KLB.
type GeneralizedResult struct {
	UserName      string
	OwnerName     string
	RepoName      string
	BaseOwnerName string
	BaseRepoName  string
	Sha           string
	Branch        string
	Number        int

	AddedOrModifiedFiles []string
}

// IsPush is used to find out if a GeneralizedResult is from a commit
// or a pull-request.
func (gr *GeneralizedResult) IsPush() bool {
	if gr.Number == 0 {
		return true
	}
	return false
}

// parseHookPullRequest gets a githubWebhook.PullRequestPayload and checks for .yml and .yaml-files
func parseHookPullRequest(ctx context.Context, payload githubWebhook.PullRequestPayload, client *authentication.Client) (*GeneralizedResult, error) {
	var result GeneralizedResult
	if !(payload.Action == "opened" ||
		payload.Action == "edited" ||
		payload.Action == "reopened" ||
		payload.Action == "synchronized") {
		return nil, nil
	}

	// head are changes, base is "old"
	result.UserName = payload.PullRequest.User.Login
	result.OwnerName = payload.PullRequest.Head.Repo.Owner.Login
	result.RepoName = payload.PullRequest.Head.Repo.Name
	result.BaseOwnerName = payload.PullRequest.Base.Repo.Owner.Login
	result.BaseRepoName = payload.PullRequest.Base.Repo.Name
	result.Branch = payload.PullRequest.Head.Ref
	result.Sha = payload.PullRequest.Head.Sha
	result.Number = int(payload.Number)

	var options = github.ListOptions{}

	files, response, err := client.GithubClient.PullRequests.ListFiles(ctx, result.BaseOwnerName, result.BaseRepoName, result.Number, &options)
	if err != nil {
		fmt.Println("Error while getting filenames:\n", err, "\n", response)
		return nil, err
	}

	for _, file := range files {
		if strings.Contains(*file.Filename, ".yml") || strings.Contains(*file.Filename, "yaml") {
			result.AddedOrModifiedFiles = append(result.AddedOrModifiedFiles, *file.Filename)
		}
	}
	return &result, nil
}

// parseHookPush gets a github.PushPayload and parses it to a GeneralizedResult-Object.
func parseHookPush(payload githubWebhook.PushPayload, client *authentication.Client) (*GeneralizedResult, error) {
	var result = GeneralizedResult{}

	result.AddedOrModifiedFiles = lookForYamlInArray(payload.HeadCommit.Added)
	modifiedFiles := lookForYamlInArray(payload.HeadCommit.Modified)

	result.AddedOrModifiedFiles = append(result.AddedOrModifiedFiles, modifiedFiles...)

	if len(result.AddedOrModifiedFiles) == 0 {
		return nil, nil
	} else {
		commitSha := payload.HeadCommit.ID
		branchRef := payload.Ref

		result.RepoName = payload.Repository.Name
		result.OwnerName = payload.Repository.Owner.Login
		result.UserName = payload.Pusher.Name
		result.Sha = commitSha
		result.Branch = branchRef

		return &result, nil
	}
}

// lookForYamlInArray looks for .yaml or .yml-files, adds them to a string-array and returns it.
func lookForYamlInArray(filesInCommit []string) []string {
	var yamlFilenames []string
	for i := 0; i < len(filesInCommit); i++ {
		if strings.Contains(filesInCommit[i], ".yaml") ||
			strings.Contains(filesInCommit[i], ".yml") {
			yamlFilenames = append(yamlFilenames, filesInCommit[i])
		}
	}
	return yamlFilenames
}

// ParseHook checks the hook for githubWebhook.PushPayload or githubWebhook.PullRequestPayload
// and passes the payloads to the appropriate methods. It ultimately returns
// a GeneralizedResult-Object.
func ParseHook(r *http.Request, secret string, client *authentication.Client) (*GeneralizedResult, error) {
	hook, err := githubWebhook.New(githubWebhook.Options.Secret(secret))
	if err != nil {
		return nil, err
	}

	payload, err := hook.Parse(r, githubWebhook.PushEvent, githubWebhook.PullRequestEvent)
	if err != nil {
		if err == githubWebhook.ErrEventNotFound {
			// This happens if the webhook sends an event that is not push or pull-request.
			fmt.Println("This event is neither a Push nor a Pull-request.\n", err)
		}
		return nil, err
	}

	var result *GeneralizedResult

	switch p := payload.(type) {

	case githubWebhook.PushPayload:
		fmt.Println("Receiving Push-Payload.")
		result, err = parseHookPush(p, client)
		if err != nil {
			return nil, err
		}
		result.Number = 0 // 0 meaning push. If pull-request, this is a non-negative non-zero number.

	case githubWebhook.PullRequestPayload:
		fmt.Println("Receiving Pull-Request-Payload.")
		result, err = parseHookPullRequest(r.Context(), p, client)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}
