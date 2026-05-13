package model

// APIkeyTableSQL is the SQL statement for creating a table corresponding to the APIkey model
var APIkeyTableSQL = `CREATE TABLE IF NOT EXISTS apikey (
	key_hash CHARACTER(64) PRIMARY KEY,
	account_id CHARACTER(36));`

// APIkey represents a row in the "apikey" table
type APIkey struct {
	KeyHash   string `db:"key_hash"`
	AccountID string `db:"account_id"`
}
