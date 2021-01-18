//parsehook parses a github-webhook. Push-Events and Pull-Request-Events are handled.
package parsehook

import (
	"fmt"
	"net/http"
	"strings"

	"gopkg.in/go-playground/webhooks.v5/github"
)

//parseHookPullRequest gets a github.PushPayload and returns AddedFilenames, ModifiedFilenames,
//and the commitSha that are parsed from the payload.
// func parseHookPullRequest(payload github.PullRequestPayload) ([]string, []string, string) {
// 	fmt.Println("Parse Hook Pull Request method")
// 	if payload.Action == "created" || payload.Action == "updated" { //updated correct?
// 		mergeCommitSha := payload.PullRequest.MergeCommitSha
// 		fmt.Println("MergeCommitSHA:", mergeCommitSha)
// 		//authenticate
// 		//getContents (of commit)
// 		//lookforyaml

// 		//modifiedFilenames := lookForYaml(payload.PullRequest)
// 		//addedFilenames := lookForYaml(payload.HeadCommit.Added)

// 		return addedFilenames, modifiedFilenames, *mergeCommitSha
// 	} else {
// 		fmt.Println("Not a newly created pull-request. Aborting.")
// 		return nil, nil, ""
// 	}
// }

//parseHookPush gets a github.PushPayload and returns AddedFilenames, ModifiedFilenames,
//and the commitSha that are parsed from the payload.
func parseHookPush(payload github.PushPayload) ([]string, []string, string) {
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

//ParseHook checks the hook for github.PushPayload or github.PullRequestPayload
//and passes the payloads to the appropriate methods. It ultimately returns
//a list of modified files, a list of added files, and the commit-SHA.
func ParseHook(r *http.Request, secret string) ([]string, []string, string) {
	hook, _ := github.New(github.Options.Secret(secret))

	payload, err := hook.Parse(r, github.PushEvent, github.PullRequestEvent)
	if err != nil {
		if err == github.ErrEventNotFound {
			//This happens if the webhook sends an event that is not push or pull-request.
			fmt.Println("This event is neither push nor pull-request.\n", err)
		}
	}
	var added []string
	var modified []string
	var commitSha string

	switch payload.(type) {

	case github.PushPayload:
		fmt.Println("Receiving Push-Payload")
		Commits := payload.(github.PushPayload)
		added, modified, commitSha = parseHookPush(Commits)
		fmt.Printf("%+v\n", Commits)

	case github.PullRequestPayload:
		fmt.Println("Receiving Pull-Request-Payload")
		pullRequest := payload.(github.PullRequestPayload)
		//added, modified, commitSha = parseHookPullRequest(pullRequest)
		fmt.Printf("%+v\n", pullRequest)
	}
	return added, modified, commitSha
}
