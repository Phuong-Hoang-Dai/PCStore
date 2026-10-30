package model

import "errors"

const MaxLimit = 50
const UserTableName = "User"

var (
	ErrUserNameOrPasswordIncorrect    = errors.New("username or password is incorrect")
	ErrCookieMissing                  = errors.New("cookie is missing")
	ErrAuthorizationHeaderWrongFormat = errors.New("authorization header is wrong format")
	ErrMissingRole                    = errors.New("missing role")
	ErrNotAllowToAccess               = errors.New("not allow to access")
	ErrCannotReadRefreshToken         = errors.New("can't read refresh_token")
	ErrEmailAlreadyRegistered         = errors.New("an account with this email already exists")
	ErrNameRequired                   = errors.New("name is required")
	ErrInvalidEmail                   = errors.New("email is invalid")
)
