package orm

import (
    "log"
    "os"
    "strings"
    "sync"

    "github.com/adhocore/urlsh/model"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var conn *gorm.DB

func pgConnect() *gorm.DB {
    dsn := os.Getenv("APP_DB_DSN")
    if dsn == "" {
        log.Fatal("Database configuration DSN missing, Pass in APP_DB_DSN env")
    }

    logLevel, env := logger.Warn, os.Getenv("APP_ENV")
    if env == "prod" {
        logLevel = logger.Silent
    } else if env == "test" {
        dsn = strings.Replace(dsn, "dbname=", "dbname=test_", 1)
    }

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logLevel),
    })

    if err != nil {
        log.Fatalf("database error: %v", err)
    }

    _ = db.AutoMigrate(&model.Keyword{}, &model.Url{})

    return db
}

var once sync.Once

func Connection() *gorm.DB {
    once.Do(func() {
        conn = pgConnect()
    })

    return conn
}
