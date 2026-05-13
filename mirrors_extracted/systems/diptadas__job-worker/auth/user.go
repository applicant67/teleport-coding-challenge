package auth

import (
	"fmt"
)

const (
	PermissionReadWrite = "READ_WRITE"
	PermissionReadOnly  = "READ_ONLY"
)

// User defines the properties of a user for using the job-worker service.
type User struct {
	UserName    string   // unique username
	Permissions []string // permissions associated with a user
}

// userDB contains list of User. It mocks a database for the POC.
// TODO: integrate with a database.
var userDB = []User{
	{
		UserName:    "alice",
		Permissions: []string{PermissionReadOnly, PermissionReadWrite},
	},
	{
		UserName:    "bob",
		Permissions: []string{PermissionReadOnly},
	},
}

// FindUser finds and returns a User object from the mock DB for a given username.
// It returns error when username do not exists.
func FindUser(username string) (User, error) {
	for _, user := range userDB {
		if user.UserName == username {
			return user, nil
		}
	}
	return User{}, fmt.Errorf("user %v not found", username)
}

// VerifyPermission tests a permission for a User.
// It returns error if permission is not granted for any role of the user.
func VerifyPermission(user User, permission string) error {
	for _, p := range user.Permissions {
		if p == permission {
			return nil
		}
	}
	return fmt.Errorf("user %v do not have permission %v", user.UserName, permission)
}

// VerifyPermissionForName tests a permission for a username.
// It returns error if username do not exists or permission is not granted for any role of the user.
func VerifyPermissionForName(username string, permission string) error {
	user, err := FindUser(username)
	if err != nil {
		return err
	}
	return VerifyPermission(user, permission)
}
