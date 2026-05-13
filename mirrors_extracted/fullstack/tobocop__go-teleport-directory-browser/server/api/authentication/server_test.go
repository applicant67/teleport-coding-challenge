package authentication

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tobocop/go-teleport-directory-browser/api/session"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var validUsername = "valid-user"
var validPassword = "valid-pass"

func TestAuthenticationHandler(t *testing.T) {
	tests := []struct {
		name string
		username string
		password string
		httpStatus int
	}{
		{"valid credentials", validUsername, validPassword, http.StatusNoContent},
		{"invalid credentials", validUsername, "invalid|pass", http.StatusUnauthorized},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("reutrns status code %d for %s", test.httpStatus, test.name), func(t *testing.T) {
			server := newMockServer(nil, "", nil)

			req := buildRequest(t, test.username, test.password)
			res := httptest.NewRecorder()

			server.AuthHandler(res, req)

			if res.Code != test.httpStatus {
				t.Errorf("got status %d but wanted %d", res.Code, test.httpStatus)
			}
		})
	}

	t.Run("starts a session and sets session cookie when authentication successful", func(t *testing.T) {
		server := newMockServer(nil, "some-session-id", nil)

		req := buildRequest(t, validUsername, validPassword)
		res := httptest.NewRecorder()

		server.AuthHandler(res, req)

		cookieValues := []string {
			"Strict",
			"HttpOnly",
			"Secure",
			"some-session-id",
			session.CookieName,
		}

		for _, val := range cookieValues {
			cookie := res.Header().Get("Set-Cookie")
			if !strings.Contains(cookie, val) {
				t.Errorf("expected cookie to contain %s, but it did not, cookie was: %v", val, cookie)
			}
		}
	})


	t.Run("returns auth header when authentication fails", func(t *testing.T) {
		server := newMockServer(nil, "", nil)

		req := buildRequest(t, validUsername, "invalid|pass")
		res := httptest.NewRecorder()

		server.AuthHandler(res, req)

		if res.Header().Get("WWW-Authenticate") == "" {
			t.Error("expected to receive WWW-Authenticate header but it was not present")
		}
	})

	t.Run("returns server error on authenticator errors", func(t *testing.T) {
		server := newMockServer(errors.New("some-error"), "", nil)

		req := buildRequest(t, validUsername, validPassword)
		res := httptest.NewRecorder()

		server.AuthHandler(res, req)

		if res.Code != http.StatusInternalServerError {
			t.Errorf("got status %d but wanted %d", res.Code, http.StatusInternalServerError)
		}
	})

	t.Run("returns server error on new session errors", func(t *testing.T) {
		server := newMockServer(nil, "", errors.New("session-error"))

		req := buildRequest(t, validUsername, validPassword)
		res := httptest.NewRecorder()

		server.AuthHandler(res, req)

		if res.Code != http.StatusInternalServerError {
			t.Errorf("got status %d but wanted %d", res.Code, http.StatusInternalServerError)
		}
	})

	t.Run("returns method not allowed when request is not a post", func(t *testing.T) {
		server := newMockServer(nil, "", nil)
		for _, method := range []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodConnect,
			http.MethodOptions,
			http.MethodTrace,
		} {
			req := httptest.NewRequest(method, "/api/authenticate", nil)
			res := httptest.NewRecorder()

			server.AuthHandler(res, req)

			if res.Code != http.StatusMethodNotAllowed {
				t.Errorf("method %v should not be allowed and is", method)
			}
		}
	})

	t.Run("returns a bad request error when the body cannot be decoded", func(t *testing.T) {
		server := newMockServer(nil, "", nil)
		req := httptest.NewRequest(http.MethodPost, "/api/authenticate", nil)
		res := httptest.NewRecorder()

		server.AuthHandler(res, req)

		if res.Code != http.StatusBadRequest {
			t.Errorf("got status %d but wanted %d", res.Code, http.StatusBadRequest)
		}
	})
}

type MockAuthenticator struct{
	err error
}

func (m *MockAuthenticator) Authenticate(username string, password string) (bool, error) {
	return username == validUsername && password == validPassword, m.err
}

type MockSessionManager struct {
	sessionId string
	err error
}

func (m *MockSessionManager) NewSession() (string, error) {
	return m.sessionId, m.err
}

func newMockServer(err error, sessionId string, sessionErr error) Server {
	return Server{
		Authenticator: &MockAuthenticator{err},
		SessionManager: &MockSessionManager{sessionId, sessionErr},
	}
}

func buildRequest(t *testing.T, username, password string) *http.Request {
	t.Helper()
	authReq := AuthRequest{username, password}
	payload, err := json.Marshal(authReq)

	if err != nil {
		t.Fatal("Json marshalling failed")
	}

	req := httptest.NewRequest(http.MethodPost, "/api/authenticate", bytes.NewBuffer(payload))
	return req
}
