// Package server contains the hook-receiving server.
package server

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/aMMokschaf/KubeLinterBot/internal/authentication"
	"github.com/aMMokschaf/KubeLinterBot/internal/config"
	"github.com/aMMokschaf/KubeLinterBot/internal/engine"
)

// SetupServer sets up the http-server.
func SetupServer(logger *log.Logger, cfg config.Config) *http.Server {
	return &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.Bot.Port),
		Handler: newServer(cfg, logWith(logger)),

		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

// Server object.
type Server struct {
	logger *log.Logger
	cfg    config.Config
	client *authentication.Client
	engine *engine.AnalysisEngine
}

// newServer creates the Server-Object.
func newServer(cfg config.Config, options ...Option) *Server {
	s := &Server{logger: log.New(ioutil.Discard, "", 0)}

	for _, o := range options {
		o(s)
	}

	s.cfg = cfg

	s.client = authentication.CreateClient(s.cfg.User.AccessToken)

	s.engine = engine.GetEngine()
	s.engine.SetClient(s.client)

	return s
}

// Option
type Option func(*Server)

// logWith creates the logger needed for the http-server.
func logWith(logger *log.Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}

// ServeHTTP waits for a github-webhook and then passes this hook to the AnalysisEngine of KubeLinterBot.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := s.engine.Analyse(r, s.cfg.User.Secret)
	if err != nil {
		s.log("Something went wrong:\n", err)
	}
}

// log logs messages
func (s *Server) log(format string, v ...interface{}) {
	s.logger.Printf(format+"\n", v...)
}
