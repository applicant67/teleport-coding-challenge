package handlers

import (
	"log"
	"net/http"

	"github.com/ibeckermayer/teleport-interview/backend/internal/auth"
	"github.com/ibeckermayer/teleport-interview/backend/internal/util"
)

// LogoutHandler handles calls to "/api/logout". Implements HandlerWithSession
type LogoutHandler struct {
	sm *auth.SessionManager
}

// NewLogoutHandler creates a new LogoutHandler
func NewLogoutHandler(sm *auth.SessionManager) *LogoutHandler {
	return &LogoutHandler{sm}
}

func (lh *LogoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := lh.sm.FromContext(r.Context())
	if err != nil {
		log.Println(err)
		util.ErrorJSON(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	// Couldn't find session. Log that this happened and return a 204
	// to alert the client to behave as if the user is logged out.
	if !lh.sm.DeleteSession(session.SessionID) {
		log.Println("logout attempted but could not find session")
	}

	w.WriteHeader(http.StatusNoContent)
}
