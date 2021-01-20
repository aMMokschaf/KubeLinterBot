//Package parsehook parses a github-webhook. Push-Events and Pull-Request-Events are handled.
package parsehook

import (
	"fmt"
	"main/internal/authentication"
	"net/http"
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	githubWebhook "gopkg.in/go-playground/webhooks.v5/github"
)

//diff datei laden, parsen: .yaml oder .yml? wenn geändert: branch name? letzter commit?

// parseHookPullRequest gets a github.PushPayload and returns AddedFilenames, ModifiedFilenames,
// and the commitSha that are parsed from the payload.
// func parseHookPullRequest(payload github.PullRequestPayload) ([]string, []string, string) {
// 	fmt.Println("Parse Hook Pull Request method")
// 	if payload.Action == "created" || payload.Action == "updated" { //updated correct?
// 		mergeCommitSha := payload.PullRequest.MergeCommitSha
// 		fmt.Println("MergeCommitSHA:", mergeCommitSha)

// 		//authenticate
// 		// githubClient := authentication.GetGithubClient()
// 		// githubClient.Repositories.GetContents(oauth2.NoContext, "aMMokschaf", "yamls", "")

// 		// var options = github.RepositoryContentGetOptions{}
// 		// _, folder, _, err := githubClient.Repositories.GetContents(oauth2.NoContext,
// 		// 	"aMMokschaf",
// 		// 	"yamls",
// 		// 	"",
// 		// 	&options)
// 		// if err != nil {
// 		// 	fmt.Println("Blablabla")
// 		// }

// 		//getContents (of commit)
// 		//githubClient.Repositories.getContents()
// 		//lookforyaml

// 		//modifiedFilenames := lookForYaml(payload.PullRequest)
// 		//addedFilenames := lookForYaml(payload.HeadCommit.Added)

// 		//return addedFilenames, modifiedFilenames, *mergeCommitSha
// 		return nil, nil, *mergeCommitSha
// 	} else {
// 		fmt.Println("Not a newly created pull-request. Aborting.")
// 		return nil, nil, ""
// 	}
// }

//parseHookPullRequest gets a githubWebhook.PullRequestPayload and checks for .yml and .yaml-files
//return should be changed
func parseHookPullRequest(payload githubWebhook.PullRequestPayload) ([]string, []string, string) {
	fmt.Println("Parse Hook Pull Request method")
	if payload.Action == "opened" ||
		payload.Action == "edited" ||
		payload.Action == "reopened" {

		var headCommitSha *string = &payload.PullRequest.Head.Sha
		githubClient := authentication.GetGithubClient()
		var options = github.ListOptions{}

		//commitFiles, response, err := githubClient.PullRequests.ListFiles(oauth2.NoContext, "aMMokschaf", "yamls", 17, &options)
		commitFiles, _, _ := githubClient.PullRequests.ListFiles(oauth2.NoContext, "aMMokschaf", "yamls", 17, &options)
		//fmt.Println("Commitfiles:", commitFiles, "\nresponse", response, "\nerr", err)
		for _, file := range commitFiles {
			if strings.Contains(*file.Filename, ".yml") || strings.Contains(*file.Filename, "yaml") {
				fmt.Println("Found yamls in PullRequest.")
				return nil, nil, *headCommitSha
			}
		}
	}
	return nil, nil, ""
}

//parseHookPush gets a github.PushPayload and returns AddedFilenames, ModifiedFilenames,
//and the commitSha that are parsed from the payload.
func parseHookPush(payload githubWebhook.PushPayload) ([]string, []string, string) {
	modifiedFilenames := lookForYaml(payload.HeadCommit.Modified)
	addedFilenames := lookForYaml(payload.HeadCommit.Added)
	commitSha := payload.HeadCommit.ID

	return addedFilenames, modifiedFilenames, commitSha
}

//lookForYaml looks for .yaml or .yml-files, adds them to a string-array and returns it.
func lookForYaml(filesInCommit []string) []string {
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
func ParseHook(r *http.Request, secret string) ([]string, []string, string) {
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

	switch payload.(type) {

	case githubWebhook.PushPayload:
		fmt.Println("Receiving Push-Payload")
		Commits := payload.(githubWebhook.PushPayload)
		added, modified, commitSha = parseHookPush(Commits)
		//fmt.Printf("%+v\n", Commits)

	case githubWebhook.PullRequestPayload:
		fmt.Println("Receiving Pull-Request-Payload")
		pullRequest := payload.(githubWebhook.PullRequestPayload)
		added, modified, commitSha = parseHookPullRequest(pullRequest)
		//fmt.Printf("%+v\n", pullRequest)
	}
	return added, modified, commitSha
}
