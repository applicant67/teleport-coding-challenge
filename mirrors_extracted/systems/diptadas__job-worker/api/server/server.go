package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"job-worker/auth"
	"net/http"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Server has router instance and other configuration to start the API server.
type Server struct {
	Port       int
	CaCert     string
	ServerCert string
	ServerKey  string
	router     *mux.Router
}

// InitializeAndRun initializes and run the API server on it's router.
func (s *Server) InitializeAndRun() {
	s.router = mux.NewRouter()
	s.setRouters()

	addr := fmt.Sprintf(":%v", s.Port)

	cert, err := ioutil.ReadFile(s.CaCert)
	if err != nil {
		log.Fatalf("could not read CA certificate, reason: %v", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert)

	httpServer := http.Server{
		Addr:    addr,
		Handler: s.router,
		TLSConfig: &tls.Config{
			ClientAuth: tls.RequireAndVerifyClientCert,
			ClientCAs:  caCertPool,
			MinVersion: tls.VersionTLS12,
		},
	}

	log.Infoln("Server started at", addr)
	log.Fatalln(httpServer.ListenAndServeTLS(s.ServerCert, s.ServerKey))
}

// setRouters sets the all required routers.
// It associates each path with http method, required permission, and request handler.
func (s *Server) setRouters() {
	s.Route("/create", http.MethodPost, auth.PermissionReadWrite, CreateJob)
	s.Route("/stop/{id}", http.MethodPut, auth.PermissionReadWrite, StopJob)
	s.Route("/status/{id}", http.MethodGet, auth.PermissionReadOnly, GetJobStatus)
}

// Route wraps the router for a HTTP method.
func (s *Server) Route(path string, method string, permission string, handler RequestHandlerFunction) {
	s.router.HandleFunc(path, s.handleRequestWithAuth(handler, permission)).Methods(method)
}

// RequestHandlerFunction defines a common type for request handler.
type RequestHandlerFunction func(w http.ResponseWriter, r *http.Request)

// handleRequestWithAuth wraps the handler function to verify user authorization before handling the request.
func (s *Server) handleRequestWithAuth(reqHandlerFunc RequestHandlerFunction, permission string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if verifyClientPermission(w, r, permission) {
			reqHandlerFunc(w, r)
		}
	}
}
