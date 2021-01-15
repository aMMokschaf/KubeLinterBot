package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"main/internal/authentication"
	"main/internal/callkubelinter"
	"main/internal/getcommit"
	"main/internal/handleresult"
	"main/internal/parsehook"
	"main/internal/postcomment"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

type config struct {
	Repository struct {
		RepoName string `yaml:"reponame"`
		User     struct {
			Username    string `yaml:"username"`
			AccessToken string `yaml:"accessToken"`
		}
	}
	Bot struct {
		Port int `yaml:"port"`
	}
}

var cfg config

//optionParser reads a config-file named "kube-linter-bot-configuration.yaml", that has
//to be located in the same folder as kube-linter-bot and parses its contents to a struct.
//A sample file is located in /samples/
func optionParser() config {
	dat, err := ioutil.ReadFile("kube-linter-bot-configuration.yaml")
	if err != nil {
		panic(err)
	}
	yaml.Unmarshal([]byte(dat), &cfg)
	//fmt.Println("Read configuration-file:\n", string(dat))
	return cfg
}

//writeOptionsToFile saves changes to the configuration to kube-linter-bot-configuration.yaml.
func writeOptionsToFile() bool {
	status := false

	//fmt.Println(cfg)
	d, err := yaml.Marshal(cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	err = ioutil.WriteFile("./kube-linter-bot-configuration.yaml", d, 0666) //TODO: Check permissions
	//fmt.Printf("%s", d)
	if err != nil {
		panic(err)
	} else {
		//fmt.Println("Setting status to true")
		status = true
	}
	return status
}

//main sets up a logger, a webHookServer, prints the address and port, starts the server
func main() {
	cfg = optionParser()
	var wg sync.WaitGroup
	if cfg.Repository.User.AccessToken == "empty" {
		wg.Add(1)
		go authentication.RunAuth(&wg)
		wg.Wait()
		cfg.Repository.User.AccessToken = authentication.GetToken()
		status := writeOptionsToFile()
		if status == false {
			fmt.Println("not written")
		}
		if status == true {
			fmt.Println("written")
		}
		//TODO: implement check if token is actually valid
	}
	logger := log.New(os.Stdout, "", 0)
	webHookServ := setupServer(logger, cfg.Bot.Port)
	logger.Printf("KubeLinterBot is listening on http://localhost%s\n", webHookServ.Addr) //TODO: Address
	webHookServ.ListenAndServe()
}

//setupServer sets up the http-server.
func setupServer(logger *log.Logger, port int) *http.Server {
	return &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: newServer(logWith(logger)),

		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

type Server struct {
	mux    *http.ServeMux
	logger *log.Logger
}

//newServer creates the http-server.
func newServer(options ...Option) *Server {
	s := &Server{logger: log.New(ioutil.Discard, "", 0)}

	for _, o := range options {
		o(s)
	}

	s.mux = http.NewServeMux()
	s.mux.HandleFunc("/", s.index)
	return s
}

type Option func(*Server)

//logWith creates the logger needed for the http-server.
func logWith(logger *log.Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}

//ServeHTTP waits for a github-webhook and then blabla TODO
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var added []string
	var modified []string
	var commitSha string
	var token string = authentication.ExtractTokenStringFromJSONToken(cfg.Repository.User.AccessToken)

	var userName = cfg.Repository.User.Username
	var repoName = cfg.Repository.RepoName

	added, modified, commitSha = parsehook.ParseHook(r)

	getcommit.GetCommit(repoName, added, modified, token)

	var klResult, klExitCode = callkubelinter.Callkubelinter()
	var hRStatus = handleresult.HandleResult(klExitCode)
	if hRStatus == 1 {
		postcomment.PostComment(token, userName, repoName, commitSha, klResult)
	}
}

//log logs messages
func (s *Server) log(format string, v ...interface{}) {
	s.logger.Printf(format+"\n", v...)
}

//TODO: Do i need this?
func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("KubeLinterBot is running here."))
}
