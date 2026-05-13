package server

import (
	"encoding/json"
	"errors"
	"job-worker/auth"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// respondJSON makes the response with payload as json format.
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, errWrite := w.Write([]byte(err.Error())); errWrite != nil {
			log.Errorln(errWrite)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, errWrite := w.Write(response); errWrite != nil {
		log.Errorln(errWrite)
	}
}

// respondMessage attaches a message with the response payload as json format.
func respondMessage(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"message": message})
}

// respondError makes the error response with payload as json format.
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}

// verifyClientPermission extracts username from the common name of the client certificate.
// It verifies the required permission for the username and returns true if permission is granted.
// If permission is not granted or, user is not found it writes the error to response body and returns false.
func verifyClientPermission(w http.ResponseWriter, r *http.Request, permission string) bool {
	username, err := getCommonName(r)
	if err != nil {
		log.Infoln(err)
		respondError(w, http.StatusUnauthorized, err.Error())
		return false
	}

	log.Infoln("Username:", username)

	user, err := auth.FindUser(username)
	if err != nil {
		log.Infoln(err)
		respondError(w, http.StatusUnauthorized, err.Error())
		return false
	}

	if err = auth.VerifyPermission(user, permission); err != nil {
		log.Infoln(err)
		respondError(w, http.StatusForbidden, err.Error())
		return false
	}

	log.Infof("User %v verified for permission %v\n", username, permission)
	return true
}

// getCommonName returns the common name of the client certificate from a http request
func getCommonName(r *http.Request) (string, error) {
	if r.TLS != nil && len(r.TLS.VerifiedChains) > 0 && len(r.TLS.VerifiedChains[0]) > 0 {
		return r.TLS.VerifiedChains[0][0].Subject.CommonName, nil
	} else {
		return "", errors.New("could not extract common name")
	}
}
