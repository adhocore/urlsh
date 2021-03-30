package cache

import (
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"

	"github.com/adhocore/urlsh/model"
	"github.com/gomodule/redigo/redis"
)

var once sync.Once
var pool *redis.Pool
var prefix = "url:"

func connect() {
	dsn := os.Getenv("REDIS_URL")
	if dsn == "" {
		return
	}

	parse, _ := url.Parse(dsn)
	user := parse.User.Username()
	pass, _ := parse.User.Password()

	if user == "h" {
		user = ""
	}

	pool = &redis.Pool{
		MaxIdle:     12,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", parse.Host, redis.DialUsername(user), redis.DialPassword(pass))
		},
	}
}

// Connection connects to redis once and returns the connection
func Connection() redis.Conn {
	once.Do(connect)

	if nil != pool {
		return pool.Get()
	}

	return nil
}

// LookupURL looks up if certain short code is popular enough to be in cache
// It returns model.URL so the request can be served right way without db hit.
func LookupURL(shortCode string) (model.URL, int) {
	var urlModel model.URL

	conn := Connection()
	if nil == conn {
		return urlModel, 0
	}

	defer conn.Close()
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

	return urlModel, http.StatusMovedPermanently
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
	conn := Connection()
	if nil == conn || (!force && hasURL(urlModel)) {
		return
	}

	defer conn.Close()
	_, _ = conn.Do("SET", urlKey(urlModel), urlValue(urlModel))
}

func hasURL(urlModel model.URL) bool {
	conn := Connection()
	if nil == conn {
		return false
	}

	defer conn.Close()
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
