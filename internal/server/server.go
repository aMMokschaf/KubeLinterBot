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
	"time"

	"main/internal/config"
	//"main/internal/engine"
	"main/internal/authentication"

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

//ServeHTTP waits for a github-webhook and then blabla TODO
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
			if len(allRepos) >= 1001 {
				break
			}
			options.Page = resp.NextPage
		}
	}
	for i, repo := range allRepos {
		fmt.Println(i, "repo fullname", repo.Owner)
	}

	// fmt.Print("\n\n\n\n\n")

	// listoptions = &github.ListOptions{Page: 2, PerPage: 13}
	// options = &github.SearchOptions{Sort: "created", Order: "asc", ListOptions: *listoptions}
	// result, resp, err = client.GithubClient.Search.Repositories(context.Background(), query, options)
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println("result:", result)
	// 	fmt.Println("resp:", resp)
	// }

	//result in generalizedresult umbauen

	//getcommit aufrufen

	//ergebnis von getcommit an callkubelinter

	//ergebnis von callkubelinter behandeln

	//CSV definieren
	columns := []string{
		"reponame", "ownername", "error1", "error2",
	}
	writer := csv.NewWriter(os.Stdout) //change to file
	if err := writer.Write(columns); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}
	writer.Flush()

	//behandeltes ergebnis in CSV schreiben
	if err := writer.Write(columns); err != nil {
		log.Fatalln("error writing record to csv:", err)
	}
	writer.Flush()

	if err := writer.Error(); err != nil {
		log.Fatal(err)
	}
}
