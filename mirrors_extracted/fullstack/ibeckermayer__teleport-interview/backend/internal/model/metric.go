package model

import "time"

// MetricTableSQL is the SQL statement for createing a table corresponding to the Metric model
var MetricTableSQL = `CREATE TABLE IF NOT EXISTS metric (
	metric_id CHARACTER(36) PRIMARY KEY,
	account_id CHARACTER(36),
	user_id CHARACTER(36),
	timestamp DATETIME
);`

// Metric represents a row in the "metric" table
type Metric struct {
	MetricID  string    `db:"metric_id"`
	AccountID string    `db:"account_id"`
	UserID    string    `db:"user_id"`
	Timestamp time.Time `db:"timestamp"`
}
