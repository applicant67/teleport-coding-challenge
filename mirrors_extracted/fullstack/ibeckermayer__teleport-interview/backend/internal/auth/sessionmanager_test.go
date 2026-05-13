package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ibeckermayer/teleport-interview/backend/internal/model"
)

func initTestSessionManager(timestr string) (*SessionManager, Session, error) {
	to, _ := time.ParseDuration(timestr)
	sm := NewSessionManager(to)
	acct := model.Account{
		AccountID:    "accountID",
		Plan:         "",
		Email:        "",
		PasswordHash: "",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now()}
	sess, err := sm.CreateSession(acct)
	return sm, sess, err
}

func TestValidSession(t *testing.T) {
	sm, sess, err := initTestSessionManager("12h")
	if err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/test", sm.WithSessionAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		return
	})))

	ts := httptest.NewTLSServer(mux)
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL+"/test", nil)
	req.Header.Set("Authorization", "Bearer "+string(sess.SessionID))
	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected %v but got %v", http.StatusOK, resp.StatusCode)
	}
}

func TestFromContext(t *testing.T) {
	sm, sess, err := initTestSessionManager("12h")
	if err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/test", sm.WithSessionAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessFromContext, err := sm.FromContext(r.Context())
		if err != nil {
			t.Fatal(err)
		}
		if !(sessFromContext.SessionID == sess.SessionID && sessFromContext.Account.AccountID == sess.Account.AccountID) {
			t.Fatal("sessFromContext.SessionID == sess.SessionID && sessFromContext.Account.AccountID == sess.Account.AccountID failed")
		}
		return
	})))

	ts := httptest.NewTLSServer(mux)
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL+"/test", nil)
	req.Header.Set("Authorization", "Bearer "+string(sess.SessionID))
	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected %v but got %v", http.StatusOK, resp.StatusCode)
	}
}

func TestUpdateSession(t *testing.T) {
	sm, sess, err := initTestSessionManager("12h")
	if err != nil {
		t.Fatal(err)
	}

	sess.Account.AccountID = "newAccountID"
	sm.UpdateSession(sess)

	checkSess, err := sm.getSession(sess.SessionID)
	if err != nil {
		t.Fatal(err)
	}
	if checkSess.Account.AccountID != sess.Account.AccountID {
		t.Fatal("checkSess.Account.AccountID != sess.Account.AccountID failed")
	}
}

func TestTimeout(t *testing.T) {
	// Set super short timeout
	to, _ := time.ParseDuration("1ns")
	sm := NewSessionManager(to)

	acct := model.Account{
		AccountID:    "accountID",
		Plan:         "",
		Email:        "",
		PasswordHash: "",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now()}

	sess, err := sm.CreateSession(acct)

	mux := http.NewServeMux()
	mux.Handle("/test", sm.WithSessionAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		return
	})))

	ts := httptest.NewTLSServer(mux)
	defer ts.Close()

	time.Sleep(1 * time.Nanosecond) // sleep to allow for timeout to happen
	req, _ := http.NewRequest("GET", ts.URL+"/test", nil)
	req.Header.Set("Authorization", "Bearer "+string(sess.SessionID))
	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected %v but got %v", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestWrongSessionID(t *testing.T) {
	to, _ := time.ParseDuration("12h")
	sm := NewSessionManager(to)

	acct := model.Account{
		AccountID:    "accountID",
		Plan:         "",
		Email:        "",
		PasswordHash: "",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now()}

	_, err := sm.CreateSession(acct)
	if err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/test", sm.WithSessionAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		return
	})))

	ts := httptest.NewTLSServer(mux)
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL+"/test", nil)
	req.Header.Set("Authorization", "Bearer "+"WrongID")
	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected %v but got %v", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestBadlyFormattedAuthHeader(t *testing.T) {
	sm, sess, err := initTestSessionManager("12h")
	if err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/test", sm.WithSessionAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		return
	})))

	ts := httptest.NewTLSServer(mux)
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL+"/test", nil)
	req.Header.Set("Authorization", string(sess.SessionID))
	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected %v but got %v", http.StatusUnauthorized, resp.StatusCode)
	}
}

func TestNoAuthHeader(t *testing.T) {
	to, _ := time.ParseDuration("12h")
	sm := NewSessionManager(to)

	acct := model.Account{
		AccountID:    "accountID",
		Plan:         "",
		Email:        "",
		PasswordHash: "",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now()}

	_, err := sm.CreateSession(acct)
	if err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/test", sm.WithSessionAuth(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		return
	})))

	ts := httptest.NewTLSServer(mux)
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL+"/test", nil)
	resp, err := ts.Client().Do(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected %v but got %v", http.StatusUnauthorized, resp.StatusCode)
	}
}
