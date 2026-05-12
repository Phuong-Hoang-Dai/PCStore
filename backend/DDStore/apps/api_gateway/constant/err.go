package constant

import "errors"

var (
	ErrUserNameOrPasswordIncorrect       = errors.New("username or password is incorrect")
	ErrCookieMissing                     = errors.New("cookie is missing")
	ErrAuthorizationHeaderWrongFormat    = errors.New("authorization header is wrong format")
	ErrMissingRole                       = errors.New("missing role")
	ErrMissingId                         = errors.New("missing Id")
	ErrNotAllowToAccess                  = errors.New("not allow to access")
	ErrCannotReadRefreshToken            = errors.New("can't read refresh_token")
	ErrCannotReadBody                    = errors.New("can't read body")
	ErrTryToParseEmptyStringToTime       = errors.New("try to parse empty string to time")
	ErrTrytoParseWrongFormatStringToTime = errors.New("try to parse wrong format string to time")
)
