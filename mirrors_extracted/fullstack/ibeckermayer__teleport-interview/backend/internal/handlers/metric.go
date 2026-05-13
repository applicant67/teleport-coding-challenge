package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/ibeckermayer/teleport-interview/backend/internal/auth"
	"github.com/ibeckermayer/teleport-interview/backend/internal/database"
	"github.com/ibeckermayer/teleport-interview/backend/internal/model"
	"github.com/ibeckermayer/teleport-interview/backend/internal/util"
)

// MetricsPostHandler handles POST calls to "api/metrics"
type MetricsPostHandler struct {
	db *database.Database
}

// NewMetricsPostHandler creates a new MetricsPostHandler
func NewMetricsPostHandler(db *database.Database) *MetricsPostHandler {
	return &MetricsPostHandler{db}
}

type metricsPostRequestBody struct {
	AccountID string    `json:"account_id"`
	UserID    string    `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
}

// Handles "/api/metrics" POST requests. Should be wrapped with WithAPIkeyAuth middlewear
func (mph *MetricsPostHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var body metricsPostRequestBody

	// Decode json
	err := util.DecodeJSONBody(w, r, &body)
	if err != nil {
		util.HandleJSONdecodeError(w, err)
		return
	}

	// Save the metric to the database
	// TODO: should metric be saved regardless of whether CreateUser below fails?
	if err := mph.db.CreateMetric(body.AccountID, body.UserID, body.Timestamp); err != nil {
		log.Println(err)
		util.ErrorJSON(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Check if this user exists
	_, err = mph.db.GetUser(body.UserID)
	if err != nil {
		// If user DNE, save new user to the database
		if err := mph.db.CreateUser(body.UserID, body.AccountID); err != nil {
			log.Println(err)
			util.ErrorJSON(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}
}

// MetricsGetHandler handles GET calls to "api/metrics"
type MetricsGetHandler struct {
	sm *auth.SessionManager
	db *database.Database
}

// NewMetricsGetHandler creates a new MetricsGetHandler
func NewMetricsGetHandler(sm *auth.SessionManager, db *database.Database) *MetricsGetHandler {
	return &MetricsGetHandler{sm, db}
}

type metricsGetResponseBody struct {
	Plan       model.Plan `json:"plan"`
	MaxUsers   int        `json:"maxUsers"`
	TotalUsers int        `json:"totalUsers"`
}

// Handles "api/metrics" GET requests. Should be wrapped with WithSessionAuth and WithAPIHeaders
func (mgh *MetricsGetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := mgh.sm.FromContext(r.Context())
	if err != nil {
		log.Println(err)
		util.ErrorJSON(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	totalUsers, err := mgh.db.CountUsers(session.Account.AccountID)
	if err != nil {
		log.Println(err)
		util.ErrorJSON(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	respBody := metricsGetResponseBody{
		session.Account.Plan,
		model.PlanMaxUsers[session.Account.Plan],
		totalUsers,
	}

	if err := json.NewEncoder(w).Encode(respBody); err != nil {
		log.Println(err)
		return
	}
}
