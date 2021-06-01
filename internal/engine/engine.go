package engine

import (
	"fmt"
	"net/http"

	"github.com/aMMokschaf/KubeLinterBot/internal/authentication"
	"github.com/aMMokschaf/KubeLinterBot/internal/callkubelinter"
	"github.com/aMMokschaf/KubeLinterBot/internal/config"
	"github.com/aMMokschaf/KubeLinterBot/internal/getcommit"
	"github.com/aMMokschaf/KubeLinterBot/internal/handleresult"
	"github.com/aMMokschaf/KubeLinterBot/internal/parsehook"
)

//AnalysisEngine
type AnalysisEngine struct{}

//GetEngine returns the AnalysisEngine-Object.
func GetEngine() *AnalysisEngine {
	var ae AnalysisEngine
	return &ae
}

//Analyse starts the processing of the payload of an incoming webhook.
func (ae *AnalysisEngine) Analyse(r *http.Request, cfg config.Config) error {
	token := cfg.User.AccessToken
	client := authentication.CreateClient(token) // use single client for "fetcher"
	result, err := parsehook.ParseHook(r, cfg.User.Secret, client)
	if err != nil {
		fmt.Println("Error while parsing hook:\n", err)
		return err
	}
	if result == nil {
		fmt.Println("Hook is of no interest to KubeLinterBot.\nKubeLinterBot is listening for Webhooks...")
		return nil
	}
	dir, err := getcommit.GetCommit(result, *client)
	if err != nil {
		return err
	}

	lintResult, exitCode := callkubelinter.CallKubelinter(dir)
	handleresult.Handle(result, lintResult, exitCode, dir, client)
	return nil
}
