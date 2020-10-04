# urlsh

**urlsh** is URL shortener application built on [Go](https://golang.org) language.

It does not use external libraries except the [`gorm`](http://gorm.io) for
[`postgres`](https://github.com/go-gorm/postgres) database.

It registers itself as Go module `github.com/adhocore/urlsh`
(however it has not been submitted to Go package registry for public usage).

## Getting source

```sh
git clone git@github.com:adhocore/urlsh.git
cd urlsh
```

## Configuring

It should be configured using env variables.

Please check [.env.example](./.env.example) for available variables and explanation.

`APP_DB_DSN` is always required and is string of the following form:

```
APP_DB_DSN=host=<str> port=<int> dbname=<str> user=<str> password=<pwd>
```

When running *urlsh* with docker-compose, the **preferred** way, `APP_DB_DSN` is
automatically set from [`POSTGRES_*`](https://hub.docker.com/_/postgres) variables.

> Please note that `urlsh` does not ship with `.env` loader so to run it in bare metal,
one needs to use `export KEY=VALUE` or `source .env` manually.

## Setting up docker

To set up dockerized `urlsh`, run the commands below:

```sh
# first time only
cp .example.env .env

# change auth token for admin if you want in `.env` file
# APP_ADMIN_TOKEN=<something crypto secure random hash>

docker-compose up
```

After a few seconds, you should be able to browse to [localhost:1000](http://localhost:1000).

## Testing

For running tests,

```sh
docker-compose exec urlsh sh -c "APP_ENV=test go test ./..."
```

---
## API Endpoints

### GET /

Index route for health-check &/or default landing.

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
    "short_url": "http://localhost:1000/qaFxz"
}
```

---
### GET /{{shortCode}}

Redirects the shortcode to original long URL.

##### Response payload

In case short code exists it responds with 302 redirect.

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
