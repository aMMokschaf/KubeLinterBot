package parsehook

import (
	"fmt"
	"net/http"
	"strings"

	"gopkg.in/go-playground/webhooks.v5/github"
)

func parseHookPullRequestToStruct(pullRequest github.PullRequestPayload) []string {
	filenames := make([]string, 2)
	filenames[0] = "Test0"
	filenames[1] = "Test1"
	//ggf.type switch: Pull request hook or push hook?
	//modified/added/deleted auf yaml prüfen
	//Dateinamen zurückgeben
	fmt.Println("Parse Hook Pull Request method")
	return filenames
}

func parseHookPush(payload github.PushPayload) ([]string, []string) {
	fmt.Println("Parse Hook Push method")
	modifiedFilenames := lookForYaml(payload.HeadCommit.Modified)
	AddedFilenames := lookForYaml(payload.HeadCommit.Added)
	return modifiedFilenames, AddedFilenames
}

func lookForYaml(filenames []string) []string {
	var modifiedFilenames []string
	for i := 0; i < len(filenames); i++ {
		if strings.Contains(filenames[i], ".yaml") ||
			strings.Contains(filenames[i], ".yml") {
			modifiedFilenames = append(modifiedFilenames, filenames[i])
		}
	}
	fmt.Println("Filenames", filenames)
	fmt.Println("Modified Filenames:", modifiedFilenames)
	return modifiedFilenames
}

func ParseHook(r *http.Request) ([]string, []string) {
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

	switch payload.(type) {

	case github.PushPayload:
		fmt.Println("push payload")
		Commits := payload.(github.PushPayload)
		added, modified = parseHookPush(Commits)
		fmt.Printf("%+v", Commits)

	case github.PullRequestPayload:
		fmt.Println("pull request payload")
		pullRequest := payload.(github.PullRequestPayload)
		//added, modified = parseHookPullRequestToStruct(pullRequest)
		fmt.Printf("%+v", pullRequest)
	}
	return added, modified
}
