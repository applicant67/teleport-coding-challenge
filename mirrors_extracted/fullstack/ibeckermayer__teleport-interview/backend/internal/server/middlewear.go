package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/ibeckermayer/teleport-interview/backend/internal/auth"
	"github.com/ibeckermayer/teleport-interview/backend/internal/util"
)

var (
	errBadAccountID = errors.New("account ID improperly formatted")
)

func getAPIkey(r *http.Request) (auth.Key, error) {
	s, err := auth.GetBearerToken(r)
	return auth.Key(s), err
}

// getAccountIDfromBody attempts to get the "account_id" field from the *http.Request body without
// mutating the *http.Request (see https://stackoverflow.com/a/47295689/6277051). Handles http error
// responses and returns an error if "account_id" field doesn't exist or other error occurs. Caller
// should return immediately on non-nil error
func getAccountIDfromBody(w http.ResponseWriter, r *http.Request) (string, error) {
	// Check that request body exists
	if r.Body == nil {
		err := errors.New("request body was nil")
		log.Println(err)
		util.ErrorJSON(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return "", err
	}

	r.Body = http.MaxBytesReader(w, r.Body, util.BodyMaxSize) // limit request body size

	// Read request body into buffer
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		util.ErrorJSON(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return "", err
	}

	// A map container to decode the JSON body into
	b := make(map[string]json.RawMessage)

	// Unmarshal json
	if err := json.Unmarshal(bodyBytes, &b); err != nil {
		log.Println(err)
		util.ErrorJSON(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return "", err
	}

	// Extract account_id
	accountIDraw, ok := b["account_id"]
	if !ok {
		err := errors.New("request body missing required field \"account_id\"")
		log.Println(err)
		util.ErrorJSON(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return "", err
	}
	accountID := string(accountIDraw)
	accountID = strings.TrimPrefix(accountID, "\"")
	accountID = strings.TrimSuffix(accountID, "\"")

	// Quick and dirty check that accountID is 36 character string.
	// TODO: change this to a regex for the uuid format
	if len(accountID) != 36 {
		log.Println(errBadAccountID)
		util.ErrorJSON(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return "", errBadAccountID
	}

	// Now that accountID is extracted, restore request body to its original state so that next can use it
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	return accountID, nil
}

// WithAPIkeyAuth is a middlewear function for protecting handlers for routes that require an API key.
// APIkey protected requests require that the sender send an API key in the Authorization header, as well
// as its corresponding account_id field in the request's body.
func (srv *Server) WithAPIkeyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the api key from the Authorization header
		key, err := getAPIkey(r)
		if err != nil {
			// Could not get APIkey, return 401
			log.Println(err)
			util.ErrorJSON(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		// Get accoutID from the request body
		accountID, err := getAccountIDfromBody(w, r)
		if err != nil {
			// getAccountIDfromBody takes care of error handling and logging, return immediately
			return
		}

		// Get corresponding row from the apikey table in the database
		apikey, err := srv.db.GetAPIkey(accountID)
		if err != nil {
			log.Printf("could not find apikey for account_id=%v", accountID)
			// w.WriteHeader(http.StatusForbidden)
			util.ErrorJSON(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		// Check that key hashes match
		if !auth.CheckKeyHash(key, apikey.KeyHash) {
			log.Printf("recieved invalid api key for account %v", accountID)
			// w.WriteHeader(http.StatusForbidden)
			util.ErrorJSON(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			return
		}

		// Request authorized, call next
		next.ServeHTTP(w, r)
	})
}

// WithAPIHeaders adds security headers to the wrapped handler
func WithAPIHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Add("Strict-Transport-Security", "max-age=631138519; includeSubDomains")
		w.Header().Add("X-Content-Type-Options", "nosniff")
		w.Header().Add("X-Permitted-Cross-Domain-Policies", "none")
		next.ServeHTTP(w, r)
	})
}

// WithHTMLHeaders adds security headers for serving an html page / spa app.
// Based on https://github.com/github/secure_headers
func WithHTMLHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: need to figure out Content-Type header that will work with X-Content-Type-Options: nosniff
		// TODO: This Content-Security-Policy may block some css that we need (untested)
		w.Header().Add("Content-Security-Policy", "default-src 'self' https:; font-src 'self' https: data:; img-src 'self' https: data:; object-src 'none'; script-src https:; style-src 'self' https: 'unsafe-inline'")
		w.Header().Add("Strict-Transport-Security", "max-age=631138519; includeSubDomains")
		w.Header().Add("X-Download-Options", "noopen")
		w.Header().Add("X-Frame-Options", "sameorigin")
		w.Header().Add("X-Permitted-Cross-Domain-Policies", "none")
		w.Header().Add("X-Xss-Protection", "1; mode=block")
		next.ServeHTTP(w, r)
	})
}
