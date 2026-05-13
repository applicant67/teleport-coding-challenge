package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/ibeckermayer/teleport-interview/backend/internal/auth"
	"github.com/ibeckermayer/teleport-interview/backend/internal/database"
	"github.com/ibeckermayer/teleport-interview/backend/internal/handlers"
)

// Config is the top level config object.
type Config struct {
	Port           int           // -port; default 8000
	CertFilePath   string        // -cert; default "../certs/localhost.crt"
	KeyFilePath    string        // -key ; default "../certs/localhost.key"
	SessionTimeout time.Duration // -sesh; default 12h
	Env            string        // -env; default "prod"
}

// Server object initializes route handlers and external connections, and serves application
type Server struct {
	cfg    Config
	router *mux.Router
	sm     *auth.SessionManager
	db     *database.Database
}

// New initializes routes and handlers and returns a ready-to-run server
func New(cfg Config) (*Server, error) {
	dbcfg := database.Config{Env: cfg.Env}
	db, err := database.New(dbcfg)
	if err != nil {
		return &Server{}, err
	}
	srv := &Server{cfg, mux.NewRouter(), auth.NewSessionManager(cfg.SessionTimeout), db}

	loginHandler := WithAPIHeaders(handlers.NewLoginHandler(srv.sm, srv.db))
	srv.router.Handle("/api/login", loginHandler).Methods("POST")

	logoutHandler := WithAPIHeaders(srv.sm.WithSessionAuth(handlers.NewLogoutHandler(srv.sm)))
	srv.router.Handle("/api/logout", logoutHandler).Methods("DELETE")

	metricsPostHandler := WithAPIHeaders(srv.WithAPIkeyAuth(handlers.NewMetricsPostHandler(srv.db)))
	srv.router.Handle("/api/metrics", metricsPostHandler).Methods("POST")

	metricsGetHandler := WithAPIHeaders(srv.sm.WithSessionAuth(handlers.NewMetricsGetHandler(srv.sm, srv.db)))
	srv.router.Handle("/api/metrics", metricsGetHandler).Methods("GET")

	upgradeHandler := WithAPIHeaders(srv.sm.WithSessionAuth(handlers.NewUpgradeHandler(srv.sm, srv.db)))
	srv.router.Handle("/api/upgrade", upgradeHandler).Methods("PATCH")

	// NOTE: It's important that this handler be registered after the other handlers, or else
	// all routes return a 404 (at least in development). TODO: figure out why this is the case.
	spaHandler := WithHTMLHeaders(handlers.NewSpaHandler("../frontend", "index.html"))
	srv.router.PathPrefix("/").Handler(spaHandler)

	return srv, nil
}

// Run starts the server
func (srv *Server) Run() error {
	log.Printf("Server listening on port %v", srv.cfg.Port)
	return http.ListenAndServeTLS(fmt.Sprintf(":%v", srv.cfg.Port), srv.cfg.CertFilePath, srv.cfg.KeyFilePath, srv.router)
}
