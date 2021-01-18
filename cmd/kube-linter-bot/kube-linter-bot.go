//Package main reads config files, and contains the hook-receiving server.
package main

import (
	"fmt"
	"log"
	"main/internal/authentication"
	"os"
	"sync"
)

//main sets up a logger, a webHookServer, prints the address and port, starts the server
func main() {
	cfg = optionParser()
	var wg sync.WaitGroup
	//TODO: implement check if token is actually valid, not just "empty"
	if cfg.Repository.User.AccessToken == "empty" {
		wg.Add(1)
		go authentication.RunAuth(&wg)
		wg.Wait()
		cfg.Repository.User.AccessToken = authentication.GetFullToken()
		status := writeOptionsToFile()
		if status == false {
			fmt.Println("Could not update configuration.")
		}
		if status == true {
			fmt.Println("Configuration updated.")
		}
	}
	logger := log.New(os.Stdout, "", 0)
	webHookServ := setupServer(logger, cfg.Bot.Port)
	logger.Printf("KubeLinterBot is listening on http://localhost%s\n", webHookServ.Addr) //TODO: Address
	webHookServ.ListenAndServe()
}
