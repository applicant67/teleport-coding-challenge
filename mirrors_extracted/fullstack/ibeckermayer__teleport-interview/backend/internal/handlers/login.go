package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ibeckermayer/teleport-interview/backend/internal/auth"
	"github.com/ibeckermayer/teleport-interview/backend/internal/database"
	"github.com/ibeckermayer/teleport-interview/backend/internal/util"
)

// LoginHandler handles calls to "/api/login". Implements http.Handler
type LoginHandler struct {
	sm *auth.SessionManager
	db *database.Database
}

// NewLoginHandler creates a new LoginHandler
func NewLoginHandler(sm *auth.SessionManager, db *database.Database) *LoginHandler {
	return &LoginHandler{sm, db}
}

type loginRequestBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginResponseBody struct {
	SessionID auth.SessionID `json:"sessionID"`
}

// Handles user login
func (lh *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body loginRequestBody

	err := util.DecodeJSONBody(w, r, &body)
	if err != nil {
		util.HandleJSONdecodeError(w, err)
		return
	}

	account, err := lh.db.GetAccountByEmail(body.Email)

	// Handle errors from attempting to retrieve the account from the database
	if err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			// No record with the given email address exists
			util.ErrorJSON(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		util.ErrorJSON(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Account retrieved, check password
	if !auth.CheckPasswordHash(body.Password, account.PasswordHash) {
		// Invalid password, unauthorized
		util.ErrorJSON(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	// Valid password, create new session
	session, err := lh.sm.CreateSession(account)
	if err != nil {
		log.Println(err)
		util.ErrorJSON(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(loginResponseBody{session.SessionID}); err != nil {
		log.Println(err)
		return
	}
}
