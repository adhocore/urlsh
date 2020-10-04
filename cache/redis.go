package cache

import (
    "log"
    "net/http"
    "os"
    "sync"

    "github.com/adhocore/urlsh/model"
    "github.com/gomodule/redigo/redis"
)

var once sync.Once
var conn redis.Conn
var prefix = "url:"

func connect() redis.Conn {
    cacheHost := os.Getenv("APP_CACHE_HOST")
    if cacheHost == "" {
        return  nil
    }

    c, err := redis.Dial("tcp", cacheHost)
    if err != nil {
        log.Printf("%v", err)
        return nil
    }

    return c
}

func Connection() redis.Conn {
    once.Do(func() {
        conn = connect()
    })

    return conn
}

func LookupURL(shortCode string) (model.URL, int) {
    var urlModel model.URL
    if nil == Connection() {
        return urlModel, 0
    }

    line, err := conn.Do("GET", urlKey(model.URL{ShortCode: shortCode}))
    if err != nil || line == nil {
        return urlModel, 0
    }

    data := string(line.([]uint8))

    // 0 = Inactive, 1 = Active
    if data[0:1] == "0" {
        return urlModel, http.StatusGone
    }

    urlModel.OriginURL = data[1:]
    urlModel.ShortCode = shortCode

    return urlModel, http.StatusFound
}

func DeactivateUrl(urlModel model.URL) {
    cacheModel, status := LookupURL(urlModel.ShortCode)
    if status == 0 {
        return
    }

    urlModel.OriginURL = cacheModel.OriginURL
    SavePopularUrl(urlModel, true)
}

func SavePopularUrl(urlModel model.URL, force bool) {
    if nil == Connection() || (!force && hasUrl(urlModel)) {
        return
    }

    _, _ = conn.Do("SET", urlKey(urlModel), urlValue(urlModel))
}

func hasUrl(urlModel model.URL) bool {
    exist, err := conn.Do("EXISTS", urlKey(urlModel))

    return err == nil && exist.(int64) > 0
}

func urlKey(urlModel model.URL) string {
    return prefix + urlModel.ShortCode
}

func urlValue(urlModel model.URL) string {
    active := "0"
    if urlModel.IsActive() {
        active = "1"
    }

    return active + urlModel.OriginURL
}
