package common

import "errors"

// Common errors
var (
    ErrInvalidUrl      = errors.New("url is invalid")
    ErrInvalidUrlLen   = errors.New("url is too short or too long, should be 7-2048 chars")
    ErrBlacklistedUrl  = errors.New("url matches blacklist pattern")
    ErrKeywordsCount   = errors.New("url must not have more than 10 keywords")
    ErrKeywordLength   = errors.New("keyword must contain 2-25 characters")
    ErrInvalidDate     = errors.New("expires_on should be in 'yyyy-mm-dd hh:mm:ss' format")
    ErrPastExpiration  = errors.New("expires_on can not be date in past")
    ErrUrlAlreadyShort = errors.New("the url is already shortened")
    ErrNoMatchingData  = errors.New("no data matching given criteria found")
    ErrTokenRequired   = errors.New("auth token is required")
    ErrTokenInvalid    = errors.New("auth token is invalid")
    ErrNoShortCode     = errors.New("the given short code is not found")
)
