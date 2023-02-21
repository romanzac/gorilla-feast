package repository

import (
	"github.com/romanzac/gorilla-feast/domain/model"
	"github.com/romanzac/gorilla-feast/middleware"
)

// UserRepository interface for basic operations with Users.
type UserRepository interface {
	Find(acct, fullname, sortQuery string, limit, offset int, noDetail bool) ([]model.User, error)
	Create(acct, fullname, pwd string) error
	Update(acct, fullname, pwd string) error
	Delete(acct string) error
	Validate(acct, pwd string) (middleware.JWTToken, error)
}
