package common

import "errors"

// Common errors
var (
	ErrInvalidURL      = errors.New("url is invalid")
	ErrInvalidURLLen   = errors.New("url is too short or too long, should be 7-2048 chars")
	ErrBlacklistedURL  = errors.New("url matches blacklist pattern")
	ErrKeywordsCount   = errors.New("keywords must not be more than 10")
	ErrKeywordLength   = errors.New("keyword must contain 2-25 characters")
	ErrInvalidDate     = errors.New("expires_on should be in 'yyyy-mm-dd hh:mm:ss' format")
	ErrPastExpiration  = errors.New("expires_on can not be date in past")
	ErrURLAlreadyShort = errors.New("url is already shortened")
	ErrNoMatchingData  = errors.New("no data matching given criteria")
	ErrTokenRequired   = errors.New("auth token is required")
	ErrTokenInvalid    = errors.New("auth token is invalid")
	ErrShortCodeEmpty  = errors.New("short_code must not be empty")
	ErrNoShortCode     = errors.New("short_code is not found")
)
