package url

import (
	"net/http"

	"github.com/adhocore/urlsh/cache"
	"github.com/adhocore/urlsh/common"
	"github.com/adhocore/urlsh/model"
	"github.com/adhocore/urlsh/orm"
	"github.com/adhocore/urlsh/request"
)

// ListURLsFilteredFromRequest gets list of urls filtered using http.Request query params
// It returns list of matching model.URL arrays and error if nothing matched.
func ListURLsFilteredFromRequest(req *http.Request) ([]model.URL, error) {
	_ = req.ParseForm()

	filter := request.URLFilter{
		ShortCode: req.Form.Get("short_code"),
		Keyword:   req.Form.Get("keyword"),
		Page:      req.Form.Get("page"),
	}

	return ListURLsFiltered(filter)
}

// ListURLsFiltered gets list of urls filtered using request.URLFilter
// It returns list of matching model.URL arrays and error if nothing matched.
func ListURLsFiltered(filter request.URLFilter) ([]model.URL, error) {
	var urls []model.URL

	limit, conn := 50, orm.Connection().Select("short_code, origin_url, hits, deleted, expires_on")
	if filter.ShortCode != "" {
		limit = 1
		conn = conn.Where("short_code = ?", filter.ShortCode)
	}

	if filter.Keyword != "" {
		conn = conn.
			Joins("LEFT JOIN url_keywords ON url_keywords.url_id = urls.id").
			Joins("LEFT JOIN keywords ON url_keywords.keyword_id = keywords.id").
			Where("keyword = ?", filter.Keyword)
	}

	if conn.Offset(filter.GetOffset(limit)).Limit(limit).Find(&urls); len(urls) == 0 {
		return urls, common.ErrNoMatchingData
	}

	return urls, nil
}

// DeleteURLFromRequest deletes url using short code from request
// It returns error on failure.
func DeleteURLFromRequest(req *http.Request) error {
	_ = req.ParseForm()

	shortCode := req.Form.Get("short_code")

	return DeleteURLByShortCode(shortCode)
}

// DeleteURLByShortCode deletes url using short code
// It returns error on failure.
func DeleteURLByShortCode(shortCode string) error {
	if shortCode == "" {
		return common.ErrShortCodeEmpty
	}

	urlModel := model.URL{ShortCode: shortCode, Deleted: true}
	result := orm.Connection().Model(urlModel).
		Where("short_code = ? AND deleted = ?", shortCode, false).
		Updates(urlModel)

	if result.RowsAffected == 0 {
		return common.ErrNoShortCode
	}

	go cache.DeactivateURL(urlModel)

	return nil
}
