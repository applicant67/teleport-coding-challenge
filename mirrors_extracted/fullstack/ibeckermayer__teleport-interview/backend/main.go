package main

import (
	"flag"
	"log"
	"time"

	"github.com/ibeckermayer/teleport-interview/backend/internal/server"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	port := flag.Int("port", 8000, "Port to serve the app from")
	certFilePath := flag.String("cert", "../certs/localhost.crt", "Relative path to a valid SSL cert")
	keyFilePath := flag.String("key", "../certs/localhost.key", "Relative path to the cert's private key")
	sessionTimeout := flag.String("sesh", "12h", "A parseable duration string (https://golang.org/pkg/time/#ParseDuration) specifying the absolute timeout value for user sessions")
	env := flag.String("env", "prod", "System environment, can be one of \"dev\" or \"prod\". The env value will determine whether the production or development database is created/used; if \"dev\", the app will seed the database with sample data for manual testing.")
	flag.Parse()

	timeout, err := time.ParseDuration(*sessionTimeout)
	if err != nil {
		log.Fatalf("failed to parse duration string for command line flag sesh=%v; see https://golang.org/pkg/time/#ParseDuration", *sessionTimeout)
	}

	cfg := server.Config{
		Port:           *port,
		CertFilePath:   *certFilePath,
		KeyFilePath:    *keyFilePath,
		SessionTimeout: timeout,
		Env:            *env}
	srv, err := server.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(srv.Run())
}
