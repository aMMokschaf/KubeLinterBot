package engine

import (
	"fmt"
	"main/internal/authentication"
	"main/internal/callkubelinter"
	"main/internal/config"
	"main/internal/getcommit"
	"main/internal/handleresult"
	"main/internal/parsehook"
	"net/http"
)

type analysisEngine struct{}

//GetEngine returns the analysisEngine-Object.
func GetEngine() *analysisEngine {
	var ae analysisEngine
	return &ae
}

//Analyse starts the processing of the payload of an incoming webhook.
func (ae *analysisEngine) Analyse(r *http.Request, cfg config.Config) error {
	var commitSha string
	var token string = cfg.User.AccessToken
	client := authentication.CreateClient(token)

	result, err := parsehook.ParseHook(r, cfg.Repositories[0].Webhook.Secret, client)
	if err != nil {
		fmt.Println("Error while parsing hook:\n", err)
		return err
	}
	if result == nil {
		fmt.Println("Hook is of no interest to KubeLinterBot.\n KubeLinterBot is listening for Webhooks...")
		return nil
	} else {
		commitSha = result.Sha
		getcommit.GetCommit(result, *client)

		var lintResult, exitCode = callkubelinter.CallKubelinter()
		handleresult.Handle(result, lintResult, exitCode, commitSha, client)
	}
	return nil
}
