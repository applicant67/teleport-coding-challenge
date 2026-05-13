package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ibeckermayer/teleport-interview/backend/internal/auth"
	"github.com/ibeckermayer/teleport-interview/backend/internal/database"
	"github.com/ibeckermayer/teleport-interview/backend/internal/model"
	"github.com/ibeckermayer/teleport-interview/backend/internal/util"
)

// UpgradeHandler handles calls to "api/upgrade"
type UpgradeHandler struct {
	sm *auth.SessionManager
	db *database.Database
}

// NewUpgradeHandler creates a new UpgradeHandler
func NewUpgradeHandler(sm *auth.SessionManager, db *database.Database) *UpgradeHandler {
	return &UpgradeHandler{sm, db}
}

type upgradHandlerResponseBody metricsGetResponseBody

// Handles "api/upgrade" calls. Should be wrapped with WithAPIHeaders and WithSessionAuth
func (uh *UpgradeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	session, err := uh.sm.FromContext(r.Context())
	if err != nil {
		log.Println(err)
		util.ErrorJSON(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Upgrade the account and set its excess users to active, grabbing the total number of users in the process
	totalUsers, err := uh.db.UpgradeAccount(session.Account.AccountID)
	if err != nil {
		log.Println(err)
		util.ErrorJSON(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// Update the session in the session manager
	session.Account.Plan = model.ENTERPRISE
	uh.sm.UpdateSession(session)

	// Build and send response body
	respBody := upgradHandlerResponseBody{
		session.Account.Plan,
		model.PlanMaxUsers[session.Account.Plan],
		totalUsers,
	}

	if err := json.NewEncoder(w).Encode(respBody); err != nil {
		log.Println(err)
		return
	}
}
