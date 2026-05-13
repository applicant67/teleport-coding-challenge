package authentication

import (
	"encoding/json"
	"github.com/tobocop/go-teleport-directory-browser/api/session"
	"log"
	"net/http"
	"time"
)

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *Server) AuthHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var authReq AuthRequest
	err := json.NewDecoder(req.Body).Decode(&authReq)
	if err != nil {
		log.Printf("AuthHandler json decode error: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	authenticated, err := s.Authenticator.Authenticate(authReq.Username, authReq.Password)
	if err != nil {
		log.Printf("AuthHandler authenticator error: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if authenticated {
		sessionId, err := s.SessionManager.NewSession()
		if err != nil {
			log.Printf("AuthHandler session manager error: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		cookie := http.Cookie{
			Name:     session.CookieName,
			Value:    sessionId,
			Expires:  time.Now().Add(session.ExpiresIn),
			Secure:   true,
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusNoContent)
	} else {
		w.Header().Set(
			"WWW-Authenticate",
			"API realm=Please enter a valid username and password to use this site.",
		)
		w.WriteHeader(http.StatusUnauthorized)
	}
}
