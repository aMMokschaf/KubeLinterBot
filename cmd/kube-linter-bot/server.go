//Package main reads config files, and contains the hook-receiving server.
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
	"strconv"
	"time"
)

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

//Server object.
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

//Option TODO
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
	var token string = cfg.Repository.User.AccessToken
	authentication.CreateClient(token)

	var userName = cfg.Repository.User.Username
	var repoName = cfg.Repository.RepoName

	added, modified, commitSha = parsehook.ParseHook(r, cfg.Webhook.Secret)
	if len(added) != 0 || len(modified) != 0 {
		getcommit.GetCommit(token, userName, repoName, commitSha, added, modified)

		var klResult, klExitCode = callkubelinter.Callkubelinter()
		var hRStatus = handleresult.HandleResult(klExitCode)
		if hRStatus == 1 {
			postcomment.PostComment(userName, repoName, commitSha, klResult)
		}
	} else {
		fmt.Println("No need to lint, as no .yml or .yaml were changed.")
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
