package auth

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ibeckermayer/teleport-interview/backend/internal/model"
	"github.com/ibeckermayer/teleport-interview/backend/internal/util"
)

var (
	// ErrSessionTimeout is returned when a caller attempts to access a session that's expired
	ErrSessionTimeout = errors.New("the session timed out")

	// ErrSessionDNE is returned when a caller attempts to access a session that doesn't exist
	ErrSessionDNE = errors.New("the session does not exist")
)

// Package specific contextKey type
type contextKey string

var theContextKey = contextKey("teleport-interview-auth")

// Session is an individual user's session
type Session struct {
	SessionID SessionID
	Account   model.Account
	Expires   time.Time
}

// SessionManager is an in-memory session store. The app should only ever create one
// of these and pass it around as a pointer
type SessionManager struct {
	store   map[SessionID]Session
	timeout time.Duration // absolute timeout for individual sessions
	mtx     sync.RWMutex  // mutex for store

	// contextKey is the key used to set and retrieve session data from a context.Context
	contextKey contextKey
}

// NewSessionManager creates a new *SessionManager
func NewSessionManager(timeout time.Duration) *SessionManager {
	return &SessionManager{
		store:      make(map[SessionID]Session),
		timeout:    timeout,
		contextKey: theContextKey}
}

// CreateSession creates a new session in the SessionManager's store, indexed by a
// new randomly generated SessionID, and expiring sm.timeout from the time it's created.
// It will return an error if the system's secure random number generator fails to function correctly.
func (sm *SessionManager) CreateSession(account model.Account) (Session, error) {
	sid, err := newSessionID()
	if err != nil {
		return Session{}, err
	}

	s := Session{sid, account, time.Now().Add(sm.timeout)}

	sm.mtx.Lock()
	defer sm.mtx.Unlock()
	sm.store[sid] = s

	return s, nil
}

// UpdateSession updates a session in the session manager with the new session passed in to it.
func (sm *SessionManager) UpdateSession(session Session) {
	sm.mtx.Lock()
	defer sm.mtx.Unlock()
	sm.store[session.SessionID] = session
}

// getSession gets a session by sessionID if it exists and isn't expired, otherwise
// it returns an empty Session object and a non-nil error
func (sm *SessionManager) getSession(sid SessionID) (Session, error) {
	sm.mtx.RLock()
	session, ok := sm.store[sid]
	sm.mtx.RUnlock()

	if !ok {
		return Session{}, ErrSessionDNE
	}

	if time.Now().After(session.Expires) {
		// Session expired, delete it from memory
		sm.DeleteSession(sid)
		return Session{}, ErrSessionTimeout
	}

	return session, nil
}

// FromContext gets a Session from a request context. WithSessionAuth-wrapped handlers should
// use this to access Session data from within the handler
func (sm *SessionManager) FromContext(ctx context.Context) (Session, error) {
	session, ok := ctx.Value(sm.contextKey).(Session)
	if !ok {
		return Session{}, errors.New("type assertion from context.Context value to Session failed")
	}
	return session, nil
}

// DeleteSession deletes a session from the session manager. Returns true if the session
// was found and deleted, or false if the session wasn't found
func (sm *SessionManager) DeleteSession(sid SessionID) bool {
	// Check whether the session exists and return false if it doesn't
	// NOTE from awly: All of this should happen under a single write lock.
	// If you put an RLock around the sm.store[sid] and a Lock around delete(sm.store, sid)
	// and have 2 concurrent DeleteSession calls, they might intertwine their read/write locks and return the wrong bool.
	sm.mtx.Lock()
	defer sm.mtx.Unlock()
	_, ok := sm.store[sid]
	if !ok {
		return ok
	}
	// Session does exist, delete it
	delete(sm.store, sid)

	return ok
}

var (
	errAuthHeaderNotFound     = errors.New("Authorization header expected but not found")
	errAuthHeaderNotFormatted = errors.New("Authorization header was improperly formatted")
)

// GetBearerToken is a helper function to retreive a token sent in standard "Bearer" format from a request
// (https://tools.ietf.org/html/rfc6750#page-5). If the request doesn't contain an Authorization
// header or the Authorization header is improperly formatted, getBearerToken returns "".
// Handlers generally shouldn't call this function, and should instead call getSessionID or
// getApiKey (TODO) to specify which type of token they are expecting.
func GetBearerToken(r *http.Request) (string, error) {
	reqToken := r.Header.Get("Authorization")
	if reqToken == "" {
		// Request did not contain an Authorization header
		return "", errAuthHeaderNotFound
	}

	splitToken := strings.Fields(reqToken)
	if len(splitToken) != 2 || splitToken[0] != "Bearer" {
		// Split failed, request may have been improperly formatted
		return "", errAuthHeaderNotFormatted
	}
	return splitToken[1], nil
}

func getSessionID(r *http.Request) (SessionID, error) {
	s, err := GetBearerToken(r)
	return SessionID(s), err
}

// WithSessionAuth is a middlewear function for protecting handlers for routes that
// require the user to be authenticated. If the user has an
func (sm *SessionManager) WithSessionAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID, err := getSessionID(r)
		if err != nil {
			// Could not get sessionID, return 401
			log.Println(err)
			util.ErrorJSON(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		session, err := sm.getSession(sessionID)
		if err != nil {
			// Session does not exist or timed out
			log.Println(err)
			util.ErrorJSON(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		// Valid session exists, add it to the context
		ctxWithSession := context.WithValue(r.Context(), sm.contextKey, session)

		// Update the http.Request with the new context and pass it to the next handler
		rWithSession := r.WithContext(ctxWithSession)
		next.ServeHTTP(w, rWithSession)
	})
}
