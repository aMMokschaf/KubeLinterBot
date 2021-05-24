//Package main reads config files, and contains the hook-receiving server.
package main

import (
	"fmt"
	"log"
	"main/internal/authentication"
	"main/internal/callkubelinter"
	"main/internal/config"
	"main/internal/server"
	"os"
)

//main sets up a logger, a webHookServer, prints the address and port, starts the server
func main() {
	//TODO argument for config file
	//TODO check if cfg-file exists
	cfg, err := config.OptionParser()
	if err != nil {
		fmt.Println("Could not read configuration-file. Please copy the file './samples/kube-linter-bot-configuration.yaml' to kube-linter-bots directory.")
		os.Exit(-1)
	}

	err = callkubelinter.CheckForKubeLinterBinary()
	if err != nil {
		os.Exit(-1)
	}
	//TODO: implement check if token is actually valid, not just "empty"
	if cfg.User.AccessToken == "empty" {
		authentication.RunAuth(*cfg)
		cfg, err = config.OptionParser()
		if err != nil {
			fmt.Println("Could not read configuration-file. Please copy the file './samples/kube-linter-bot-configuration.yaml' to kube-linter-bots directory.")
		}
	}
	logger := log.New(os.Stdout, "", 0)
	webHookServ := server.SetupServer(logger, *cfg)
	logger.Printf("KubeLinterBot is listening on http://localhost%s\n", webHookServ.Addr) //TODO: Address
	webHookServ.ListenAndServe()
}
