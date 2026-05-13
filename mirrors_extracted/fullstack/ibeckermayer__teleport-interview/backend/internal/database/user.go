package database

import (
	"errors"
	"time"

	"github.com/ibeckermayer/teleport-interview/backend/internal/model"
)

var (
	// ErrOrphanedUser is returned if caller attempts to create a user associated with an account_id that DNE.
	ErrOrphanedUser = errors.New("attempted to create an orphaned user")
)

// GetUser retrieves a user from the database by user_id. Returns an error if user DNE.
func (db *Database) GetUser(userID string) (model.User, error) {
	u := model.User{}
	err := db.db.Get(&u, "SELECT * FROM user WHERE user_id=$1", userID)
	return u, err
}

// CountUsers counts how many users are associated with the accountID
func (db *Database) CountUsers(accountID string) (int, error) {
	var c int
	err := db.db.Get(&c, "SELECT count(*) FROM user WHERE account_id=$1", accountID)
	return c, err
}

// CreateUser creates a new user with userID associated with accountID. Determines
// whether the new User is active based on if the associated account has reached the
// user limit on its current plan.
func (db *Database) CreateUser(userID, accountID string) error {
	createUserUpgradeAccountLock.Lock()
	defer createUserUpgradeAccountLock.Unlock()

	account, err := db.GetAccount(accountID)
	if err != nil {
		// Could not find account, don't create orphaned user
		return err
	}

	count, err := db.CountUsers(accountID)
	if err != nil {
		return err
	}

	user := &model.User{
		UserID:    userID,
		AccountID: accountID,
		IsActive:  count+1 <= model.PlanMaxUsers[account.Plan],
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return db.insertUser(user)
}

func (db *Database) insertUser(u *model.User) error {
	_, err := db.db.NamedExec("INSERT INTO user (user_id, account_id, is_active, created_at, updated_at) VALUES (:user_id, :account_id, :is_active, :created_at, :updated_at)", u)
	return err
}
