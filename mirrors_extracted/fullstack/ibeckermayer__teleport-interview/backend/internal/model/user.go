package model

import "time"

// UserTableSQL is the SQL statement for createing a table corresponding to the User model
var UserTableSQL = `CREATE TABLE IF NOT EXISTS user (
	user_id CHARACTER(36) PRIMARY KEY,
	account_id CHARACTER(36),
	is_active INTEGER,
	created_at DATETIME,
	updated_at DATETIME
);`

// User represents a row in the "user" table
type User struct {
	UserID    string    `db:"user_id"`
	AccountID string    `db:"account_id"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
