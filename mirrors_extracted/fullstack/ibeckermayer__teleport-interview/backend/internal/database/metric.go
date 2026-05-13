package database

import (
	"time"

	"github.com/pborman/uuid"

	"github.com/ibeckermayer/teleport-interview/backend/internal/model"
)

func (db *Database) insertMetric(m *model.Metric) error {
	_, err := db.db.NamedExec("INSERT INTO metric (metric_id, account_id, user_id, timestamp) VALUES (:metric_id, :account_id, :user_id, :timestamp)", m)
	return err
}

// CreateMetric adds a new metric to the "metric" table in the database
func (db *Database) CreateMetric(accountID, userID string, timestamp time.Time) error {
	metric := &model.Metric{
		MetricID:  uuid.New(),
		AccountID: accountID,
		UserID:    userID,
		Timestamp: timestamp,
	}

	return db.insertMetric(metric)
}
