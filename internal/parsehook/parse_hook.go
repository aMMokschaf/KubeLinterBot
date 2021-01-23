//Package parsehook parses a github-webhook. Push-Events and Pull-Request-Events are handled.
package parsehook

import (
	"context"
	"fmt"
	"main/internal/authentication"
	"net/http"
	"strings"

	"github.com/google/go-github/github"
	githubWebhook "gopkg.in/go-playground/webhooks.v5/github"
)

type PrSourceBranchInformation struct {
	UserName  string
	RepoName  string
	BranchRef string
	Needed    bool
}

func getPRI(payload githubWebhook.PullRequestPayload) PrSourceBranchInformation {
	var pri PrSourceBranchInformation
	pri.UserName = payload.PullRequest.User.Login
	pri.RepoName = payload.PullRequest.Head.Repo.Name
	pri.BranchRef = payload.PullRequest.Head.Ref
	if payload.PullRequest.Head.Repo.FullName != payload.PullRequest.Base.Repo.FullName {
		pri.Needed = true
	} else {
		pri.Needed = false
	}
	return pri
}

//parseHookPullRequest gets a githubWebhook.PullRequestPayload and checks for .yml and .yaml-files
//return should be changed
func parseHookPullRequest(payload githubWebhook.PullRequestPayload) ([]string, []string, string, string, PrSourceBranchInformation) {
	fmt.Println("Parse Hook Pull Request method")
	fmt.Println(payload)
	pullRequestInformation := getPRI(payload)
	// fmt.Println(payload)
	// fmt.Println(payload.PullRequest.Number)
	// fmt.Println(payload.PullRequest.Head.User.Login)
	// fmt.Println(payload.Repository.Name)
	var headCommitSha *string
	var branchRef string = payload.PullRequest.Head.Ref
	if payload.Action == "opened" ||
		payload.Action == "edited" ||
		payload.Action == "reopened" ||
		payload.Action == "synchronized" {

		headCommitSha = &payload.PullRequest.Head.Sha
		githubClient := authentication.GetGithubClient()
		var options = github.ListOptions{}

		//commitFiles, response, err := githubClient.PullRequests.ListFiles(oauth2.NoContext, "aMMokschaf", "yamls", 17, &options)
		commitFiles, response, err := githubClient.PullRequests.ListFiles(context.Background(), payload.PullRequest.Head.User.Login, payload.Repository.Name, int(payload.Number), &options)
		fmt.Println("Commitfiles:", commitFiles, "\nresponse", response, "\nerr", err)
		for _, file := range commitFiles {
			if strings.Contains(*file.Filename, ".yml") || strings.Contains(*file.Filename, "yaml") {
				fmt.Println("Found yamls in PullRequest. Return HeadCommitSha.", *headCommitSha)
				return nil, nil, *headCommitSha, branchRef, pullRequestInformation
			}
		}
	}
	var empty string = ""
	headCommitSha = &empty
	return nil, nil, *headCommitSha, branchRef, pullRequestInformation
}

//parseHookPush gets a github.PushPayload and returns AddedFilenames, ModifiedFilenames,
//and the commitSha that are parsed from the payload.
func parseHookPush(payload githubWebhook.PushPayload) ([]string, []string, string, string) {
	fmt.Println("Entering parseHookPush")
	modifiedFilenames := lookForYamlInArray(payload.HeadCommit.Modified)
	addedFilenames := lookForYamlInArray(payload.HeadCommit.Added)
	commitSha := payload.HeadCommit.ID
	branchRef := *&payload.Ref

	fmt.Println("ModifiedFiles:", modifiedFilenames)
	fmt.Println("AddedFilenames:", addedFilenames)
	fmt.Println("commitSha:", commitSha)
	fmt.Println("branchref:", branchRef)
	return addedFilenames, modifiedFilenames, commitSha, branchRef
}

//lookForYamlInArray looks for .yaml or .yml-files, adds them to a string-array and returns it.
func lookForYamlInArray(filesInCommit []string) []string {
	var yamlFilenames []string
	for i := 0; i < len(filesInCommit); i++ {
		if strings.Contains(filesInCommit[i], ".yaml") ||
			strings.Contains(filesInCommit[i], ".yml") {
			yamlFilenames = append(yamlFilenames, filesInCommit[i])
		}
	}
	//fmt.Println("Filenames", filesInCommit)
	//fmt.Println("Modified Filenames:", yamlFilenames)
	return yamlFilenames
}

//ParseHook checks the hook for githubWebhook.PushPayload or githubWebhook.PullRequestPayload
//and passes the payloads to the appropriate methods. It ultimately returns
//a list of modified files, a list of added files, and the commit-SHA.
func ParseHook(r *http.Request, secret string) ([]string, []string, string, string, PrSourceBranchInformation) {
	hook, _ := githubWebhook.New(githubWebhook.Options.Secret(secret))

	payload, err := hook.Parse(r, githubWebhook.PushEvent, githubWebhook.PullRequestEvent)
	if err != nil {
		if err == githubWebhook.ErrEventNotFound {
			//This happens if the webhook sends an event that is not push or pull-request.
			fmt.Println("This event is neither push nor pull-request.\n", err)
		}
	}
	var added []string
	var modified []string
	var commitSha string
	var branchRef string
	var pullRequestInformation PrSourceBranchInformation
	pullRequestInformation.Needed = false

	switch payload.(type) {

	case githubWebhook.PushPayload:
		fmt.Println("Receiving Push-Payload")
		Commits := payload.(githubWebhook.PushPayload)
		added, modified, commitSha, branchRef = parseHookPush(Commits)
		//fmt.Printf("%+v\n", Commits)

	case githubWebhook.PullRequestPayload:
		fmt.Println("Receiving Pull-Request-Payload")
		pullRequest := payload.(githubWebhook.PullRequestPayload)
		added, modified, commitSha, branchRef, pullRequestInformation = parseHookPullRequest(pullRequest)
		//fmt.Printf("%+v\n", pullRequest)
	}
	return added, modified, commitSha, branchRef, pullRequestInformation
}
