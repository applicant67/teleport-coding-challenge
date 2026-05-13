package util

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
)

// BodyMaxSize is the maximum size of a json body
var BodyMaxSize int64 = 1048576 // 1MB

type malformedRequest struct {
	status int    // HTTP status
	logMsg string // Detailed message for logging, not to be returned to clients
}

// Generic HTTP status message for sending to clients
func (mr malformedRequest) Error() string {
	return mr.logMsg
}

// DecodeJSONBody is a modified version of https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body
// Attempts to decode a json request body into dst, returning an error if it fails
func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "" {
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			msg := "Content-Type header is not application/json"
			return &malformedRequest{status: http.StatusUnsupportedMediaType, logMsg: msg}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, BodyMaxSize)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		return &malformedRequest{status: http.StatusBadRequest, logMsg: err.Error()}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return &malformedRequest{status: http.StatusBadRequest, logMsg: msg}
	}

	return nil
}

// HandleJSONdecodeError is a helperr function for responding to requests that cause decodeJSONBody to throw an error.
// Logs a detailed error message and responds to the client with a generic HTTP error message.
// The caller should ensure no further writes are done to w.
func HandleJSONdecodeError(w http.ResponseWriter, err error) {
	var mr *malformedRequest
	if errors.As(err, &mr) {
		log.Println(mr)
		ErrorJSON(w, http.StatusText(mr.status), mr.status)
	} else {
		log.Println(err)
		ErrorJSON(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

type errorJSON struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type errorJSONbody struct {
	Error errorJSON `json:"error"`
}

// ErrorJSON is a replacement for http.Error for use by endpoints with Content-Type header "application/json".
// http.Error should not be used because it overwrites the Content-Type header to "text/plain"
func ErrorJSON(w http.ResponseWriter, error string, code int) {
	w.WriteHeader(code)
	ejb := errorJSONbody{errorJSON{code, error}}
	if err := json.NewEncoder(w).Encode(ejb); err != nil {
		log.Println(err)
		return
	}
}
