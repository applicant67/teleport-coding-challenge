package database

import (
	"time"

	"github.com/ibeckermayer/teleport-interview/backend/internal/auth"
	"github.com/ibeckermayer/teleport-interview/backend/internal/model"
)

func (db *Database) insertAccount(a *model.Account) error {
	_, err := db.db.NamedExec("INSERT INTO account (account_id, plan, email, password_hash, created_at, updated_at) VALUES (:account_id, :plan, :email, :password_hash, :created_at, :updated_at)", a)
	return err
}

// UpgradeAccount upgrades an account from the FREE to the ENTERPRISE plan. It also updates
// any users in that account that were previously inactive to active. Returns the total number of users
// for the given accountID for ease of use by the UpgradeHandler.
// TODO: handle case when there wind up being more users than the ENTERPRISE plan allows
func (db *Database) UpgradeAccount(accountID string) (int, error) {
	createUserUpgradeAccountLock.Lock()
	defer createUserUpgradeAccountLock.Unlock()

	tx, err := db.db.Begin()
	if err != nil {
		return 0, err
	}
	// Update the account
	_, err = tx.Exec("UPDATE account SET plan=$1 WHERE account_id=$2", model.ENTERPRISE, accountID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	_, err = tx.Exec("UPDATE user SET is_active=$1 WHERE is_active=$2 AND account_id=$3", true, false, accountID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	totalUsers, err := db.CountUsers(accountID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return totalUsers, tx.Commit()
}

// GetAccount retrieves an Account from the database by accountID
func (db *Database) GetAccount(accountID string) (model.Account, error) {
	a := model.Account{}
	err := db.db.Get(&a, "SELECT * FROM account WHERE account_id=$1", accountID)
	return a, err
}

// GetAccountByEmail retrieves an Account from the database by email address
func (db *Database) GetAccountByEmail(email string) (model.Account, error) {
	a := model.Account{}
	err := db.db.Get(&a, "SELECT * FROM account WHERE email=$1", email)
	return a, err
}

// CreateAccount creates a new account and saves it in the database
func (db *Database) CreateAccount(accountID, email, password string) error {
	passwordHash, err := auth.HashPassword(password)
	if err != nil {
		return err
	}
	account := &model.Account{
		AccountID:    accountID,
		Plan:         model.FREE,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err = db.insertAccount(account)
	return err
}
