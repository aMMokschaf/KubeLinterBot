//Package server contains the hook-receiving server.
package server

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"main/internal/callkubelinter"
	"main/internal/config"
	"main/internal/getcommit"

	//"main/internal/engine"
	"main/internal/authentication"
	"main/internal/parsehook"

	"encoding/csv"

	"github.com/google/go-github/github"
)

//SetupServer sets up the http-server.
func SetupServer(logger *log.Logger, cfg config.Config) *http.Server {
	return &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.Bot.Port),
		Handler: newServer(cfg, logWith(logger)),

		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

//Server object.
type Server struct {
	mux    *http.ServeMux
	logger *log.Logger
	cfg    config.Config
}

//newServer creates the Server-Object.
func newServer(cfg config.Config, options ...Option) *Server {
	s := &Server{logger: log.New(ioutil.Discard, "", 0)}

	for _, o := range options {
		o(s)
	}

	s.mux = http.NewServeMux()
	s.mux.HandleFunc("/", s.index)
	s.mux.HandleFunc("/eval", s.eval)

	s.cfg = cfg
	return s
}

//Option
type Option func(*Server)

//logWith creates the logger needed for the http-server.
func logWith(logger *log.Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}

//ServeHTTP waits for a github-webhook and then errorString TODO
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// ae := engine.GetEngine()
	// err := ae.Analyse(r, s.cfg)
	// if err != nil {
	// 	s.log("Something went wrong:\n", err)
	// }
	// //TODO response?
	//s.eval(w, r)
	s.mux.ServeHTTP(w, r)
}

//log logs messages
func (s *Server) log(format string, v ...interface{}) {
	s.logger.Printf(format+"\n", v...)
}

//TODO: Do i need this?
func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("KubeLinterBot is running here."))
}

//eval-function to get data from github
func (s *Server) eval(w http.ResponseWriter, r *http.Request) {
	fmt.Println("eval-Funktion")
	w.Write([]byte("KubeLinterBot-evaluation is running here."))

	//client abrufen
	var token string = s.cfg.User.AccessToken
	client := authentication.CreateClient(token)

	//query-package aufrufen, query erstellen
	query := "extension:yaml OR extension:yml kubernetes OR k8s"
	//listoptions := &github.ListOptions{Page: 1, PerPage: 13}
	listoptions := &github.ListOptions{PerPage: 100}
	options := &github.SearchOptions{Sort: "created", Order: "asc", ListOptions: *listoptions}
	var allRepos []github.Repository
	for {
		result, resp, err := client.GithubClient.Search.Repositories(context.Background(), query, options)
		if err != nil {
			fmt.Println(err)
		} else {
			allRepos = append(allRepos, result.Repositories...)
			if resp.NextPage == 0 {
				break
			}
			//if len(allRepos) >= 1001 {
			if len(allRepos) >= 1 {
				break
			}
			options.Page = resp.NextPage
		}
	}

	file, err := os.Create("data.csv")
	if err != nil {
		log.Fatalln("error creating csv:", err)
	}
	defer file.Close()

	//CSV definieren
	columns := []string{
		"number",
		"reponame",
		"ownername",
		"dangling-service",
		"deprecated-service-account-field",
		"drop-net-raw-capability",
		"env-var-secret",
		"mismatching selector",
		"no-anti-affinity",
		"no-extensions-v1beta",
		"no-read-only-root-fs",
		"non-existent-service-account",
		"privileged-container",
		"run-as-non-root",
		"ssh-port",
		"unset-cpu-requirements",
		"unset-memory-requirements",
	}

	writer := csv.NewWriter(file)
	if err := writer.Write(columns); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}
	writer.Flush()

	var genRes []parsehook.GeneralizedResult
	for _, repo := range allRepos {
		var newRepo parsehook.GeneralizedResult
		newRepo.UserName = ""
		newRepo.OwnerName = *repo.Owner.Login
		newRepo.RepoName = *repo.Name
		newRepo.BaseOwnerName = ""
		newRepo.BaseRepoName = ""
		newRepo.Branch = "master" //getBranch? Irgendwie mainBranch holen
		id := strconv.Itoa(int(*repo.ID))
		newRepo.Sha = id
		newRepo.Number = 0

		_, dir, _, err := client.GithubClient.Repositories.GetContents(context.Background(), newRepo.OwnerName, newRepo.RepoName, "", nil)
		if err != nil {
			fmt.Println("getcontents newRepo", err, newRepo.OwnerName, newRepo.RepoName)
		}
		//fmt.Println(dir)
		for _, file := range dir {
			filename := file.GetPath()
			//fmt.Println("filename:", filename)
			if strings.Contains(filename, ".yaml") || strings.Contains(filename, ".yml") {
				newRepo.AddedOrModifiedFiles = append(newRepo.AddedOrModifiedFiles, filename)
				fmt.Println(filename, "true")
			}
		}
		genRes = append(genRes, newRepo)
	}
	//fmt.Println("dateinamen:", genRes[1].AddedOrModifiedFiles)
	for i, repo := range genRes {
		getcommit.GetCommit(&repo, *client)

		erg, err := callkubelinter.CallKubelinter(repo.Sha)
		if err != nil {
			fmt.Println("kubelinter err", err)
		}
		fmt.Println("ergebnis:", string(erg))

		var errorArray [14]int

		errorMsgs := strings.Split(string(erg), "\n")
		fmt.Println("errormsgs", errorMsgs[0], errorMsgs[1])
		for _, msg := range errorMsgs {
			if msg == "" {
				fmt.Println("skip rest of loop, msg")
				continue
			}
			fmt.Println("msg:", msg)

			var indexBegin int = strings.Index(msg, "check: ")
			if indexBegin == -1 {
				break
			}
			var subString string = msg[indexBegin : len(msg)-1]
			indexBegin = 7
			var indexEnd int = strings.Index(subString, ",")

			//var indexEnd int = strings.Index(msg[strings.Index(msg, ","):], ",")
			// var indexEnd int = strings.Index(msg, ",")
			// var blaString string = msg[indexEnd+1 : len(msg)-1]
			// var indexBegin int = strings.Index(blaString, "check: ")
			// indexBegin += 7 //sets index to begin of check-name
			// fmt.Println("blastring:", blaString)
			// indexEnd = strings.Index(blaString, ",")
			var errorString string
			fmt.Println("indexbegin", indexBegin, "indexEnd", indexEnd)
			if indexBegin == -1 || indexEnd == -1 || indexEnd < indexBegin {
				errorString = "none"
			} else {
				errorString = subString[indexBegin:indexEnd]
			}
			fmt.Println("errorstring:", errorString)

			switch errorString {
			case "dangling-service":
				fmt.Println(errorString)
				errorArray[0] += 1
			case "deprecated-service-account-field":
				fmt.Println(errorString)
				errorArray[1] += 1
			case "drop-net-raw-capability":
				fmt.Println(errorString)
				errorArray[2] += 1
			case "env-var-secret":
				fmt.Println(errorString)
				errorArray[3] += 1
			case "mismatching-selector":
				fmt.Println(errorString)
				errorArray[4] += 1
			case "no-anti-affinity":
				fmt.Println(errorString)
				errorArray[5] += 1
			case "no-extensons-v1beta":
				fmt.Println(errorString)
				errorArray[6] += 1
			case "no-read-only-root-fs":
				fmt.Println(errorString)
				errorArray[7] += 1
			case "non-existent-service-account":
				fmt.Println(errorString)
				errorArray[8] += 1
			case "privileged-container":
				fmt.Println(errorString)
				errorArray[9] += 1
			case "run-as-non-root":
				fmt.Println(errorString)
				errorArray[10] += 1
			case "ssh-port":
				fmt.Println(errorString)
				errorArray[11] += 1
			case "unset-cpu-requirements":
				fmt.Println(errorString)
				errorArray[12] += 1
			case "unset-memory-requirements":
				fmt.Println(errorString)
				errorArray[13] += 1
			default:
				//Wert für keine Probleme?
			}
		}
		fmt.Println(errorArray)

		dataset := []string{strconv.Itoa(i), repo.RepoName, repo.OwnerName}
		for i := range errorArray {
			number := errorArray[i]
			text := strconv.Itoa(number)
			dataset = append(dataset, text)
		}

		if err := writer.Write(dataset); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
		writer.Flush()

		if err := writer.Error(); err != nil {
			log.Fatal(err)
		}
	}
}
