package parsehook

import (
	"fmt"
	"net/http"
	"strings"

	"gopkg.in/go-playground/webhooks.v5/github"
)

//TODO: This whole method
//parseHookPullRequest gets a github.PushPayload and returns AddedFilenames, ModifiedFilenames,
//and the commitSha that are parsed from the payload.
// func parseHookPullRequest(payload github.PullRequestPayload) ([]string, []string, string) {
// 	fmt.Println("Parse Hook Pull Request method")
// 	modifiedFilenames := lookForYaml(payload.HeadCommit.Modified)
// 	addedFilenames := lookForYaml(payload.HeadCommit.Added)
// 	commitSha := payload.HeadCommit.ID

// 	return addedFilenames, modifiedFilenames, commitSha
// }

//parseHookPush gets a github.PushPayload and returns AddedFilenames, ModifiedFilenames,
//and the commitSha that are parsed from the payload.
func parseHookPush(payload github.PushPayload) ([]string, []string, string) {
	fmt.Println("Parse Hook Push method")
	modifiedFilenames := lookForYaml(payload.HeadCommit.Modified)
	addedFilenames := lookForYaml(payload.HeadCommit.Added)
	commitSha := payload.HeadCommit.ID

	return addedFilenames, modifiedFilenames, commitSha
}

//TODO: better variables
//lookForYaml looks for .yaml or .yml-files, adds them to a string-array and returns it.
func lookForYaml(filenames []string) []string {
	var modifiedFilenames []string
	for i := 0; i < len(filenames); i++ {
		if strings.Contains(filenames[i], ".yaml") ||
			strings.Contains(filenames[i], ".yml") {
			modifiedFilenames = append(modifiedFilenames, filenames[i])
		}
	}
	//fmt.Println("Filenames", filenames)
	//fmt.Println("Modified Filenames:", modifiedFilenames)
	return modifiedFilenames
}

//ParseHook checks the hook for github.PushPayload or github.PullRequestPayload
//and passes the payloads to the appropriate methods. It ultimately returns
//a list of modified files, a list of added files, and the commit-SHA.
func ParseHook(r *http.Request) ([]string, []string, string) {
	fmt.Println("parse hook method")
	hook, _ := github.New(github.Options.Secret("testsecret"))

	payload, err := hook.Parse(r, github.PushEvent, github.PullRequestEvent)
	if err != nil {
		if err == github.ErrEventNotFound {
			// ok event wasn;t one of the ones asked to be parsed
			fmt.Println(err)
		}
	}
	var added []string
	var modified []string
	var commitSha string

	switch payload.(type) {

	case github.PushPayload:
		fmt.Println("push payload")
		Commits := payload.(github.PushPayload)
		added, modified, commitSha = parseHookPush(Commits)
		fmt.Printf("%+v", Commits)

	case github.PullRequestPayload:
		fmt.Println("pull request payload")
		pullRequest := payload.(github.PullRequestPayload)
		//added, modified, commitSha = parseHookPullRequest(pullRequest)
		fmt.Printf("%+v", pullRequest)
	}
	return added, modified, commitSha
}
