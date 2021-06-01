//Package server contains the hook-receiving server.
package server

import (
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/aMMokschaf/KubeLinterBot/internal/config"
	"github.com/aMMokschaf/KubeLinterBot/internal/engine"
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
	ae := engine.GetEngine() // move ae to field of Server struct
	err := ae.Analyse(r, s.cfg)
	if err != nil {
		s.log("Something went wrong:\n", err)
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
