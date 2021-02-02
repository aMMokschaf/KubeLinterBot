//Package main reads config files, and contains the hook-receiving server.
package main

import (
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
	cfg := config.OptionParser()

	err := callkubelinter.CheckForKubeLinterBinary()
	if err != nil {
		//TODO exit server?
	}
	//var wg sync.WaitGroup
	//TODO: implement check if token is actually valid, not just "empty"
	if cfg.User.AccessToken == "empty" {
		//wg.Add(1)
		/*go*/
		authentication.RunAuth(cfg) //&wg)
		//wg.Wait()
		//authObj := authentication.CreateClient()
		cfg = config.OptionParser()
	}
	logger := log.New(os.Stdout, "", 0)
	webHookServ := server.SetupServer(logger, cfg)
	logger.Printf("KubeLinterBot is listening on http://localhost%s\n", webHookServ.Addr) //TODO: Address
	webHookServ.ListenAndServe()
}
