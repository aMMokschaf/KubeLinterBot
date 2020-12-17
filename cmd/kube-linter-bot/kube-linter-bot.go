package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"main/internal/callkubelinter"
	"main/internal/getcommit"
	"main/internal/handleresult"
	"main/internal/parsehook"
	"main/internal/postcomment"
	"net/http"
	"os"
	"strconv"
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

func optionParser() config {
	dat, err := ioutil.ReadFile("kube-linter-bot-configuration.yaml")
	if err != nil {
		panic(err)
	}
	var cfg config
	yaml.Unmarshal([]byte(dat), &cfg)
	fmt.Println("Read configuration-file:\n", string(dat))
	return cfg
}

//Sets up a logger, a webHookServer, prints the address and port, starts the server
func main() {
	var cfg = optionParser()
	fmt.Println(cfg.Bot.Port)
	logger := log.New(os.Stdout, "", 0)
	webHookServ := setupServer(logger, cfg.Bot.Port)
	logger.Printf("KubeLinterBot is listening on http://localhost%s\n", webHookServ.Addr) //TODO: Address
	webHookServ.ListenAndServe()

	/*
		react to webhook
			check if .yaml was changed
				Call KubeLinter
				if exit-code 0
					no review comment? Some other kind of feedback?
				else
					interpret Linter-output
					review comment via github
			else
				do nothing? Feedback?
	*/
}

//Setup method, needs an already set up logger and returns a http.Server-Pointer
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

func logWith(logger *log.Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	/*if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		s.log("Only POST allowed.")
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body) //TODO ReadAll may be bad for large messages
	if err != nil {
		log.Fatal(err)
	}

	s.log("Webhook received.")
	makeJSON(s, reqBody)*/

	var added []string
	var modified []string

	added, modified = parsehook.ParseHook(r)

	getcommit.GetCommit(added, modified)
	callkubelinter.Callkubelinter()
	handleresult.HandleResult()
	postcomment.PostComment()

}

/*
//TODO: Parse JSON. Marshal? Decode?
func makeJSON(s *Server, body []byte) {
	s.log("\n%s", body)
}*/

func (s *Server) log(format string, v ...interface{}) {
	s.logger.Printf(format+"\n", v...)
}

func (s *Server) index(w http.ResponseWriter, r *http.Request) {
	//w.Write([]byte("KubeLinterBot is running here."))
	//parsehook.ParseHook(r)

}
