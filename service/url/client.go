package url

import (
    "encoding/json"
    "net/http"
    "os"

    "github.com/adhocore/urlsh/common"
    "github.com/adhocore/urlsh/model"
    "github.com/adhocore/urlsh/orm"
    "github.com/adhocore/urlsh/request"
    "github.com/adhocore/urlsh/util"
    "gorm.io/gorm"
)

// CreateUrlShortCodeFromRequest creates a new short code for url given in http.Request
// It uses expires_on date and keywords from http.Request if available.
// It returns created short code or error if any.
func CreateUrlShortCodeFromRequest(req *http.Request) (string, error) {
    var input request.UrlInput

    if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
        return "", err
    }

    return CreateUrlShortCode(input)
}

// CreateUrlShortCodeFromRequest creates a new short code using request.UrlInput
// It returns created short code or error if any.
func CreateUrlShortCode(input request.UrlInput) (string, error) {
    if shortCode, err := validateUrlInput(input, allowDupeUrl()); err != nil {
        return shortCode, err
    }

    shortCode := getUniqueShortCode()
    expiresOn, _ := input.GetExpiresOn()

    orm.Connection().Create(&model.Url{
        ShortCode: shortCode,
        OriginUrl: input.Url,
        Keywords:  mapKeywords(input.Keywords),
        ExpiresOn: expiresOn,
    })

    return shortCode, nil
}

// LookupOriginUrl looks up origin url from shortCode
// It returns origin url if exists and is active, http error code otherwise.
func LookupOriginUrl(shortCode string) (string, int) {
    var urlModel model.Url

    if status := orm.Connection().Where("short_code = ?", shortCode).First(&urlModel); status.RowsAffected == 0 {
        return "", http.StatusNotFound
    }

    if !urlModel.IsActive() {
        return "", http.StatusGone
    }

    return urlModel.OriginUrl, http.StatusFound
}

// IncrementHits increments hit counter for given shortCode just before serving redirection
func IncrementHits(shortCode string) {
    var urlModel model.Url

    orm.Connection().Model(&urlModel).
        Where("short_code = ?", shortCode).
        UpdateColumn("hits", gorm.Expr("hits + ?", 1))
}

// allowDupeUrl checks is app is configured to allow dupe url
func allowDupeUrl() bool {
    return os.Getenv("APP_ALLOW_DUPE_URL") == "1"
}

// validateUrlInput validates given request.UrlInput
// It returns existing short code for input url if exists and validation error if any.
func validateUrlInput(input request.UrlInput, allowDup bool) (string, error) {
    if err := input.Validate(); err != nil || allowDup {
        return "", err
    }

    if shortCode := getShortCodeByOriginUrl(input.Url); shortCode != "" {
        return shortCode, common.ErrUrlAlreadyShort
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
    var urlModel model.Url

    status := orm.Connection().Select("short_code").Where("short_code = ?", shortCode).First(&urlModel)

    return status.RowsAffected > 0
}

// getShortCodeByOriginUrl gets short code for given origin url
// It returns short code string.
func getShortCodeByOriginUrl(originUrl string) string {
    var urlModel model.Url

    orm.Connection().Select("short_code").
        Where("origin_url = ? AND deleted = ?", originUrl, false).
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
