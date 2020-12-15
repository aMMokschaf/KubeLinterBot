package parsehook

import (
	"fmt"
	"net/http"

	"gopkg.in/go-playground/webhooks.v5/github"
)

func parseHookPullRequest(pullRequest github.PullRequestPayload) []string {
	filenames := make([]string, 2)
	filenames[0] = "Test0"
	filenames[1] = "Test1"
	//ggf.type switch: Pull request hook or push hook?
	//modified/added/deleted auf yaml prüfen
	//Dateinamen zurückgeben
	fmt.Println("Parse Hook Pull Request method")
	return filenames
}

func parseHookPush(commits github.PushPayload) []string {
	filenames := make([]string, 2)
	filenames[0] = "Test0"
	filenames[1] = "Test1"
	//ggf.type switch: Pull request hook or push hook?
	//modified/added/deleted auf yaml prüfen
	//Dateinamen zurückgeben
	fmt.Println("Parse Hook Push method")
	return filenames
}

func checkforChangedYAML() {

}

func ParseHook(r *http.Request) {
	fmt.Println("parse hook method")
	hook, _ := github.New(github.Options.Secret("testsecret"))

	payload, err := hook.Parse(r, github.ReleaseEvent, github.PullRequestEvent)
	if err != nil {
		if err == github.ErrEventNotFound {
			// ok event wasn;t one of the ones asked to be parsed
			fmt.Println(err)
		}
	}
	switch payload.(type) {

	case github.PushPayload:
		fmt.Println("push payload")
		Commits := payload.(github.PushPayload)
		//parseHookPush(Commits)
		fmt.Printf("%+v", Commits)

	case github.PullRequestPayload:
		fmt.Println("pull request payload")
		pullRequest := payload.(github.PullRequestPayload)
		//parseHookPullRequest(pullRequest)
		fmt.Printf("%+v", pullRequest)
	}
}
