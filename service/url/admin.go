package url

import (
    "net/http"

    "github.com/adhocore/urlsh/common"
    "github.com/adhocore/urlsh/model"
    "github.com/adhocore/urlsh/orm"
    "github.com/adhocore/urlsh/request"
)

func ListUrlsFilteredFromRequest(req *http.Request) ([]model.Url, error) {
    _ = req.ParseForm()

    filter := request.UrlFilter{
        ShortCode: req.Form.Get("short_code"),
        Keyword:   req.Form.Get("keyword"),
        Page:      req.Form.Get("page"),
    }

    return ListUrlsFiltered(filter)
}

func ListUrlsFiltered(filter request.UrlFilter) ([]model.Url, error) {
    var urls []model.Url

    limit, conn := 50, orm.Connection().Select("short_code, origin_url, hits, expires_on")
    if filter.ShortCode != "" {
        limit = 1
        conn  = conn.Where("short_code = ?", filter.ShortCode)
    } else {
        conn.Offset(filter.GetOffset(limit))
    }

    if filter.Keyword != "" {
        conn = conn.
            Joins("LEFT JOIN url_keywords ON url_keywords.url_id = urls.id").
            Joins("LEFT JOIN keywords ON url_keywords.keyword_id = keywords.id").
            Where("keyword = ?", filter.Keyword)
    }

    if conn.Limit(limit).Find(&urls); len(urls) == 0 {
        return urls, common.ErrNoMatchingData
    }

    return urls, nil
}
