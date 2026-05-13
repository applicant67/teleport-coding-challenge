package authentication

import (
	"github.com/tobocop/go-teleport-directory-browser/api/user"
	"golang.org/x/crypto/bcrypt"
)

type Authenticator interface {
	Authenticate(username string, password string) (bool, error)
}

type staticCredentialsAuthenticator struct {
	UserStore user.Store
}

func (sa *staticCredentialsAuthenticator) Authenticate(username string, password string) (bool, error) {
	found, storedPassword := sa.UserStore.GetUserPassword(username)
	return found && bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(password)) == nil, nil
}
