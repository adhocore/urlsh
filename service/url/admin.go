package url

import (
    "net/http"

    "github.com/adhocore/urlsh/common"
    "github.com/adhocore/urlsh/model"
    "github.com/adhocore/urlsh/orm"
    "github.com/adhocore/urlsh/request"
)

// ListUrlsFilteredFromRequest gets list of urls filtered using http.Request query params
// It returns list of matching model.Url arrays and error if nothing matched.
func ListUrlsFilteredFromRequest(req *http.Request) ([]model.Url, error) {
    _ = req.ParseForm()

    filter := request.UrlFilter{
        ShortCode: req.Form.Get("short_code"),
        Keyword:   req.Form.Get("keyword"),
        Page:      req.Form.Get("page"),
    }

    return ListUrlsFiltered(filter)
}

// ListUrlsFiltered gets list of urls filtered using request.UrlFilter
// It returns list of matching model.Url arrays and error if nothing matched.
func ListUrlsFiltered(filter request.UrlFilter) ([]model.Url, error) {
    var urls []model.Url

    limit, conn := 50, orm.Connection().Select("short_code, origin_url, hits, deleted, expires_on")
    if filter.ShortCode != "" {
        limit = 1
        conn  = conn.Where("short_code = ?", filter.ShortCode)
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

func DeleteUrlFromRequest(req *http.Request) error {
    _ = req.ParseForm()

    shortCode := req.Form.Get("short_code")

    return DeleteUrlByShortCode(shortCode)
}

func DeleteUrlByShortCode(shortCode string) error {
    result := orm.Connection().Model(model.Url{}).
        Where("short_code = ? AND deleted = ?", shortCode, false).
        Updates(model.Url{Deleted: true})

    if result.RowsAffected == 0 {
        return common.ErrNoShortCode
    }

    return nil
}
