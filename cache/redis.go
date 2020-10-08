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

    user, pass := os.Getenv("APP_CACHE_USER"), os.Getenv("APP_CACHE_PASS")
    c, err := redis.Dial("tcp", cacheHost, redis.DialUsername(user), redis.DialPassword(pass))
    if err != nil {
        log.Printf("%v", err)
        return nil
    }

    return c
}

// Connection connects to redis once and returns the connection
func Connection() redis.Conn {
    once.Do(func() {
        conn = connect()
    })

    return conn
}

// LookupURL looks up if certain short code is popular enough to be in cache
// It returns model.URL so the request can be served right way without db hit.
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

    urlModel.OriginURL = data[1:]
    urlModel.ShortCode = shortCode

    // 0 = Inactive, 1 = Active
    if data[0:1] == "0" {
        return urlModel, http.StatusGone
    }

    return urlModel, http.StatusFound
}

// DeactivateURL deactivates cache of an expired/deleted model.URL
// PS, this operation is always cached so Gone (410) can be served without db hit.
func DeactivateURL(urlModel model.URL) {
    cacheModel, _ := LookupURL(urlModel.ShortCode)

    if urlModel.OriginURL == "" {
        urlModel.OriginURL = cacheModel.OriginURL
    }

    SavePopularURL(urlModel, true)
}

// SavePopularURL saves a model.URL to cache
// If force is passed, it saves even if already exists
func SavePopularURL(urlModel model.URL, force bool) {
    if nil == Connection() || (!force && hasURL(urlModel)) {
        return
    }

    _, _ = conn.Do("SET", urlKey(urlModel), urlValue(urlModel))
}

func hasURL(urlModel model.URL) bool {
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
