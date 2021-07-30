# urlsh

[![Latest Version](https://img.shields.io/github/release/adhocore/urlsh.svg?style=flat-square)](https://github.com/adhocore/urlsh/releases)
[![Software License](https://img.shields.io/badge/license-MIT-brightgreen.svg?style=flat-square)](LICENSE)
[![lint](https://github.com/adhocore/urlsh/actions/workflows/lint-github-action.yml/badge.svg)](https://github.com/adhocore/urlsh/actions/workflows/lint-github-action.yml)
[![test](https://github.com/adhocore/urlsh/actions/workflows/test-github-action.yml/badge.svg)](https://github.com/adhocore/urlsh/actions/workflows/test-github-action.yml)
[![Donate 15](https://img.shields.io/badge/donate-paypal-blue.svg?style=flat-square&label=donate+15)](https://www.paypal.me/ji10/15usd)
[![Donate 25](https://img.shields.io/badge/donate-paypal-blue.svg?style=flat-square&label=donate+25)](https://www.paypal.me/ji10/25usd)
[![Donate 50](https://img.shields.io/badge/donate-paypal-blue.svg?style=flat-square&label=donate+50)](https://www.paypal.me/ji10/50usd)
[![Tweet](https://img.shields.io/twitter/url/http/shields.io.svg?style=social)](https://twitter.com/intent/tweet?text=URL+shortener+app+built+on+Golang&url=https://github.com/adhocore/urlsh&hashtags=go,golang,url,urlshortener)


**[urlsh](https://urlssh.xyz)** is URL shortener application built on [Go](https://golang.org) language.

It does not use external libraries except the [`gorm`](http://gorm.io) for
[`postgres`](https://github.com/go-gorm/postgres) database and minimal redis
module [redigo](https://github.com/gomodule/redigo).

![URLSH](./assets/urlsh.png)

# Usage

#### Web client

Just visit [urlssh.xyz](https://urlssh.xyz). You can also integrate API for programmatic usage.

#### Terminal client

If you live inside the command line terminal, you can install a client:
```sh
go get github.com/adhocore/urlsh/cmd/urlshc
```

Make sure `$GOPATH` or `$HOME/go/bin` is in your `$PATH` or `%path%`, then use it like so:
```sh
# show help
urlshc -h

# shorten a url
urlshc -url 'https://your.long.url/here'

# shorten a url with keywords and expiry date
urlshc -url 'https://your.long.url/here' -keywords 'word-1,word2' -expires '2022-12-31'

# if you have self hosted or local instance of urlsh, pass in URLSH_HOST (with domain only)
URLSH_HOST=https://your.urlsh-instance.com urlshc -url 'https://your.long.url/here'
```

Read below for self hosting, API integration and/or contributing to *urlsh*.

___

## Getting source

```sh
git clone git@github.com:adhocore/urlsh.git
cd urlsh
```

## Configuring

It should be configured using env variables.

Please check [.env.example](./.env.example) for available variables and explanation.

`DATABASE_URL` is always required and is string of the following form:

```
DATABASE_URL=postgres://user:pass@host:port/dbname
```

When running *urlsh* with docker-compose, the **preferred** way, `DATABASE_URL` is
automatically set from [`POSTGRES_*`](https://hub.docker.com/_/postgres) variables.

> Please note that `urlsh` does not ship with `.env` loader so to run it in bare metal,
one needs to use `export KEY=VALUE` or `source .env` manually.

If `REDIS_URL` is provided with `redis://host:port`, then it will cache popular urls,
the popularity threshold is hard coded to minimal value in `common.constant`.
The cache will sync when short url is deleted and sometimes when expired.
The hit counters will however sync in realtime.

## Setting up docker

To set up dockerized `urlsh`, run the commands below:

```sh
# first time only
cp .example.env .env

# change auth token for admin if you want in `.env` file
# APP_ADMIN_TOKEN=<something crypto secure random hash>

docker-compose up
```

After a few seconds, you should be able to browse to [localhost:2000](http://localhost:2000).

## Testing

For running tests,

```sh
docker-compose exec urlsh sh -c "APP_ENV=test go test ./..."

# for coverage
docker-compose exec urlsh sh -c "APP_ENV=test go test -cover ./..."
```

`APP_ENV=test` is not required but ensures that tests are run against clone database with name prefixed by `test_`.
Normally test db is already prepared if you use `docker-compose`.

---
## API Endpoints

### GET /status

Status route for health/status check.

#### Response payload

```json
{
    "status": 200,
    "message": "it works"
}
```

---
### POST /api/urls

Creates a new short code for given URL.

#### Request example

- url: string, required, http/https/ftp only
- expires_on: string, optional, utc date `yyyy-mm-dd hh:mm:ss`, default=`9999-01-01`
- keywords: array of strings, 2-25 chars each, max 10 keywords

```json
{
    "url": "http://somedomain.com/some/very/long/url",
    "expires_on": "",
    "keywords": ["key", "word"]
}
```

#### Response example

```json
{
    "status": 200,
    "short_code": "qaFxz",
    "short_url": "http://localhost:2000/qaFxz"
}
```

> If env var `APP_ALLOW_DUPE_URL` is set to 0 or empty, then trying to shorten same URL again
will return status 409 and payload will contain existing `short_code`.
> However if existing `short_code` is deleted, it will be shortened as usual.

---
### GET /{shortCode}

Redirects the shortcode to original long URL.

##### Response payload

In case short code exists it responds with 301 redirect.

If the short code is expired or deleted, it responds like so:

```json
{
    "status": 410,
    "message": "requested resource is not available"
}
```

---
### GET /api/admin/urls

#### Authentication

Token required in `Authorization` header like so:
```ini
Authorization: Bearer <token>
```

#### Request query

The query params are *optional*.

```ini
page=<int>
short_code=<str>
keyword=<str>
```

*Examples:*

- `/api/admin/urls?short_code=somecode`
- `/api/admin/urls?page=1&keyword=something`

#### Response example

Response contains multiple matching url object inside `urls` array.

```json
{
    "status": 200,
    "urls": [
        {
            "short_code": "X5JkFd",
            "origin_url": "http://somedomain.com/some/very/long/url",
            "hits": 1,
            "is_deleted": false,
            "expires_on": "9999-01-01T00:00:00Z"
        }
    ]
}
```

---
### DELETE /api/admin/urls

#### Authentication

Token required in `Authorization` header like so:
```
Authorization: Bearer <token>
```

#### Request query

Query param `short_code` is requied.

*Example*: `/api/admin/urls?short_code=somecode`

#### Response example

If delete success:

```json
{
    "status": 200,
    "deleted": true
}
```

If the code does not exist:

```json
{
    "status": 404,
    "message": "the given short code is not found"
}
```

---
### Using postman

**urlsh** comes with [postman](./postman) collection and environment to aid manual testing of endpoints.

Open the postman app, click `Import`  at top left, select `Folder` and drag/choose `postman` folder of this repo.
You may need to adjust the `token` in postman `urlsh` env if you have configured `APP_ADMIN_TOKEN`.

The collection comes with post/pre request hooks for requests so you can just run the endpoints one after another in postman UI.

> For `redirect` request, you have to disable postman follow redirects from `Settings > General > Automatically follow redirects`.

### License

Please check [license](./LICENSE) file.

---
### Other projects
My other golang projects you might find interesting and useful:

- [**gronx**](https://github.com/adhocore/gronx) - Lightweight, fast and dependency-free Cron expression parser (due checker), task scheduler and/or daemon for Golang (tested on v1.13 and above) and standalone usage.
- [**fast**](https://github.com/adhocore/fast) - Check your internet speed with ease and comfort right from the terminal
