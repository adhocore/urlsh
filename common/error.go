package common

import "errors"

// Common errors
var (
    ErrInvalidURL      = errors.New("url is invalid")
    ErrInvalidURLLen   = errors.New("url is too short or too long, should be 7-2048 chars")
    ErrBlacklistedURL  = errors.New("url matches blacklist pattern")
    ErrKeywordsCount   = errors.New("url must not have more than 10 keywords")
    ErrKeywordLength   = errors.New("keyword must contain 2-25 characters")
    ErrInvalidDate     = errors.New("expires_on should be in 'yyyy-mm-dd hh:mm:ss' format")
    ErrPastExpiration  = errors.New("expires_on can not be date in past")
    ErrURLAlreadyShort = errors.New("the url is already shortened")
    ErrNoMatchingData  = errors.New("no data matching given criteria found")
    ErrTokenRequired   = errors.New("auth token is required")
    ErrTokenInvalid    = errors.New("auth token is invalid")
    ErrShortCodeEmpty  = errors.New("short_code must not be empty")
    ErrNoShortCode     = errors.New("the given short_code is not found")
)
