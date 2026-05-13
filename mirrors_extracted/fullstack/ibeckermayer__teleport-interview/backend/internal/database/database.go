package database

import (
	"log"
	"os"
	"sync"

	"github.com/ibeckermayer/teleport-interview/backend/internal/auth"
	"github.com/ibeckermayer/teleport-interview/backend/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/pborman/uuid"
)

// createUserUpgradeAccountLock is used to ensure data integrity during calls of CreateUser and UpgradeAccount.
// If CreateUser gets called in the middle of the transaction in UpgradeAccount, the new user might be created with
// is_active set to false when it should really be set to true (since it's corresponding Account was being upgraded).
var createUserUpgradeAccountLock = sync.Mutex{}

// Config is a database config object.
// Env determines whether the production or development database is created/used;
// if \"dev\", the app will seed the database with sample data for manual testing.
type Config struct {
	Env string
}

// Database is a handle to the database layer
type Database struct {
	cfg Config
	db  *sqlx.DB
}

// New creates a new *Database and initializes it's schema.
// Set filldb to true to fill the database with some fake data (for development purposes).
func New(cfg Config) (*Database, error) {
	dbfile := "./teleport-interview-" + cfg.Env + ".db"
	if cfg.Env == "dev" {
		// Reset db for every dev restart
		os.Remove(dbfile)
	}
	sqlxdb, err := sqlx.Open("sqlite3", dbfile)
	if err != nil {
		return nil, err
	}

	// force a connection and test that it worked
	if err = sqlxdb.Ping(); err != nil {
		return nil, err
	}

	db := &Database{cfg: cfg, db: sqlxdb}
	if err = db.init(); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *Database) init() error {
	// Create all tables if they don't exist
	if _, err := db.db.Exec(model.AccountTableSQL); err != nil {
		return err
	}

	if _, err := db.db.Exec(model.APIkeyTableSQL); err != nil {
		return err
	}

	if _, err := db.db.Exec(model.MetricTableSQL); err != nil {
		return err
	}

	if _, err := db.db.Exec(model.UserTableSQL); err != nil {
		return err
	}

	if db.cfg.Env == "dev" {
		// Fill the db with data for development purposes
		devNamePwd := "dev@goteleport.com"
		fakeiotTestNamePwd := "test@goteleport.com"
		devAcctID := uuid.New()
		fakeiotTestAcctID := "testacct-0000-0000-0000-000000000000"
		devKey, _ := auth.NewKey()
		fakeiotTestKey, _ := auth.NewKey()

		db.CreateAccount(devAcctID, devNamePwd, devNamePwd)
		db.CreateAccount(fakeiotTestAcctID, fakeiotTestNamePwd, fakeiotTestNamePwd)
		db.CreateAPIkey(devKey, devAcctID)
		db.CreateAPIkey(fakeiotTestKey, fakeiotTestAcctID)

		log.Printf("Created dev account with account_id=%v, username/pwd=%v, and token=%v", devAcctID, devNamePwd, devKey)
		log.Printf("Created fakeiot test account with account_id=%v, username/pwd=%v, and token=%v", fakeiotTestAcctID, fakeiotTestNamePwd, fakeiotTestKey)

	}

	return nil
}
