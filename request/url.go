package request

import (
    "net/url"
    "regexp"
    "time"

    "github.com/adhocore/urlsh/common"
)

const UrlBlackListRegex = "(xxx)"

// UrlInput defines structure for create short code url request
type UrlInput struct {
    Url        string    `json:"url" binding:"required"`
    ExpiresOn  string    `json:"expires_on"`
    Keywords   []string  `json:"keywords"`
}

// Validate validates the url input before saving to db
func (input UrlInput) Validate() error {
    if l := len(input.Url); l < 7 || l > 2048 {
        return common.ErrInvalidUrlLen
    }

    if match, _ := regexp.MatchString("^(f|ht)tps?://+", input.Url); !match {
        return common.ErrInvalidUrl
    }

    if _, err := url.ParseRequestURI(input.Url); err != nil {
        return common.ErrInvalidUrl
    }

    if match, _ := regexp.MatchString(UrlBlackListRegex, input.Url); match {
        return common.ErrBlacklistedUrl
    }

    if len(input.Keywords) > 10 {
        return common.ErrKeywordsCount
    }

    for _, word := range input.Keywords {
        if l := len(word); l < 2 || l > 25 {
            return common.ErrKeywordLength
        }
    }

    return input.ValidateExpiry()
}

// ValidateExpiry validates expires_on date if not empty
func (input UrlInput) ValidateExpiry() error {
    if input.ExpiresOn == "" {
        return nil
    }

    if len(input.ExpiresOn) != len(common.DateLayout) {
        return common.ErrInvalidDate
    }

    expiresOn, err := input.GetExpiresOn()
    if err != nil {
        return common.ErrInvalidDate
    }

    if expiresOn.In(time.UTC).Before(time.Now().In(time.UTC)) {
        return common.ErrPastExpiration
    }

    return nil
}

// GetExpiresOn gets date time instance or error if parse fails
func (input UrlInput) GetExpiresOn() (time.Time, error) {
   if input.ExpiresOn == "" {
       return time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC), nil
   }

   return time.ParseInLocation(common.DateLayout, input.ExpiresOn, time.UTC)
}
