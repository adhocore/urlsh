package url

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/adhocore/urlsh/cache"
	"github.com/adhocore/urlsh/common"
	"github.com/adhocore/urlsh/model"
	"github.com/adhocore/urlsh/orm"
	"github.com/adhocore/urlsh/request"
	"github.com/adhocore/urlsh/util"
	"gorm.io/gorm"
)

// CreateURLShortCodeFromRequest creates a new short code for url given in http.Request
// It uses expires_on date and keywords from http.Request if available.
// It returns created short code or error if any.
func CreateURLShortCodeFromRequest(req *http.Request) (string, error) {
	var input request.URLInput

	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		return "", err
	}

	return CreateURLShortCode(input)
}

// CreateURLShortCode creates a new short code using request.URLInput
// It returns created short code or error if any.
func CreateURLShortCode(input request.URLInput) (string, error) {
	if shortCode, err := validateURLInput(input, allowDupeURL()); err != nil {
		return shortCode, err
	}

	shortCode := getUniqueShortCode()
	expiresOn, _ := input.GetExpiresOn()

	orm.Connection().Create(&model.URL{
		ShortCode: shortCode,
		OriginURL: input.URL,
		Keywords:  mapKeywords(input.Keywords),
		ExpiresOn: expiresOn,
	})

	return shortCode, nil
}

// LookupOriginURL looks up origin url from shortCode
// It returns origin url if exists and is active, http error code otherwise.
func LookupOriginURL(shortCode string) (model.URL, int, bool) {
	var urlModel model.URL

	if cacheModel, status := cache.LookupURL(shortCode); status != 0 {
		return cacheModel, status, true
	}

	if status := orm.Connection().Where("short_code = ?", shortCode).First(&urlModel); status.RowsAffected == 0 {
		return urlModel, http.StatusNotFound, false
	}

	if !urlModel.IsActive() {
		if !urlModel.Deleted {
			go cache.DeactivateURL(urlModel)
		}

		return urlModel, http.StatusGone, false
	}

	return urlModel, http.StatusMovedPermanently, false
}

// IncrementHits increments hit counter for given shortCode just before serving redirection
func IncrementHits(urlModel model.URL) {
	orm.Connection().Model(&urlModel).
		Where("short_code = ?", urlModel.ShortCode).
		UpdateColumn("hits", gorm.Expr("hits + ?", 1))

	if urlModel.Hits+1 >= common.PopularHits {
		cache.SavePopularURL(urlModel, false)
	}
}

// allowDupeURL checks is app is configured to allow dupe url
func allowDupeURL() bool {
	return os.Getenv("APP_ALLOW_DUPE_URL") == "1"
}

// validateURLInput validates given request.URLInput
// It returns existing short code for input url if exists and validation error if any.
func validateURLInput(input request.URLInput, allowDup bool) (string, error) {
	if err := input.Validate(); err != nil || allowDup {
		return "", err
	}

	if shortCode := getShortCodeByOriginURL(input.URL); shortCode != "" {
		return shortCode, common.ErrURLAlreadyShort
	}

	return "", nil
}

// getUniqueShortCode gets unique random string to use as short code
// It checks db to ensure it is really unique and returns short code string.
func getUniqueShortCode() string {
	shortCode := util.RandomString(common.ShortCodeLength)

	for {
		if !isExistingShortCode(shortCode) {
			return shortCode
		}

		shortCode = util.RandomString(common.ShortCodeLength)
	}
}

// isExistingShortCode checks if given short code exists in db
// It returns bool.
func isExistingShortCode(shortCode string) bool {
	var urlModel model.URL

	status := orm.Connection().Select("short_code").Where("short_code = ?", shortCode).First(&urlModel)

	return status.RowsAffected > 0
}

// getShortCodeByOriginURL gets short code for given origin url
// It returns short code string.
func getShortCodeByOriginURL(originURL string) string {
	var urlModel model.URL

	orm.Connection().Select("short_code").
		Where("origin_url = ? AND deleted = ?", originURL, false).
		First(&urlModel)

	return urlModel.ShortCode
}

// mapKeywords maps input keyword array to model arrays
// It returns array of model.Keyword.
func mapKeywords(words []string) []model.Keyword {
	var Keywords []model.Keyword

	for _, word := range words {
		var keyword model.Keyword
		orm.Connection().FirstOrCreate(&keyword, model.Keyword{Keyword: word})
		Keywords = append(Keywords, keyword)
	}

	return Keywords
}
