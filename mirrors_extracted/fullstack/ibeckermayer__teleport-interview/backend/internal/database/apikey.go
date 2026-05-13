package database

import (
	"github.com/ibeckermayer/teleport-interview/backend/internal/auth"
	"github.com/ibeckermayer/teleport-interview/backend/internal/model"
)

func (db *Database) insertAPIkey(apikey *model.APIkey) error {
	_, err := db.db.NamedExec("INSERT INTO apikey (key_hash, account_id) VALUES (:key_hash, :account_id)", apikey)
	return err
}

// GetAPIkey gets an apikey from the database by accountID
func (db *Database) GetAPIkey(accountID string) (model.APIkey, error) {
	k := model.APIkey{}
	err := db.db.Get(&k, "SELECT * FROM apikey WHERE account_id=$1", accountID)
	return k, err
}

// CreateAPIkey creates a new apikey entry
// TODO: could check that accouuntID exists in account table
func (db *Database) CreateAPIkey(key auth.Key, accountID string) error {
	keyHash := auth.HashKey(key)
	apikey := &model.APIkey{KeyHash: keyHash, AccountID: accountID}
	return db.insertAPIkey(apikey)
}
