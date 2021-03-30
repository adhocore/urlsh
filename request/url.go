package request

import (
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/adhocore/urlsh/common"
)

// URLBlackListRegex is regex to filter unwanted urls
const URLBlackListRegex = "(xxx)"

// URLInput defines structure for create short code url request
type URLInput struct {
	URL       string   `json:"url" binding:"required"`
	ExpiresOn string   `json:"expires_on"`
	Keywords  []string `json:"keywords"`
}

// URLFilter defines structure for short code list and search request
type URLFilter struct {
	ShortCode string `json:"short_code"`
	Keyword   string `json:"keyword"`
	Page      string `json:"page"`
}

// Validate validates the url input before saving to db
// It returns error if something is not valid.
func (input URLInput) Validate() error {
	if l := len(input.URL); l < 7 || l > 2048 {
		return common.ErrInvalidURLLen
	}

	if match, _ := regexp.MatchString("^(f|ht)tps?://+", input.URL); !match {
		return common.ErrInvalidURL
	}

	if _, err := url.ParseRequestURI(input.URL); err != nil {
		return common.ErrInvalidURL
	}

	if match, _ := regexp.MatchString(URLBlackListRegex, input.URL); match {
		return common.ErrBlacklistedURL
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
// It returns error if expiry date is not valid.
func (input URLInput) ValidateExpiry() error {
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
func (input URLInput) GetExpiresOn() (time.Time, error) {
	if input.ExpiresOn == "" {
		return time.Date(9999, 1, 1, 0, 0, 0, 0, time.UTC), nil
	}

	return time.ParseInLocation(common.DateLayout, input.ExpiresOn, time.UTC)
}

// GetOffset gets normalized pagination offset
func (filter URLFilter) GetOffset(limit int) int {
	if filter.Page == "" {
		return 0
	}

	page, err := strconv.Atoi(filter.Page)
	if err != nil || page < 2 {
		return 0
	}

	return (page - 1) * limit
}
