package orm

import (
    "log"
    "os"
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

    logLevel := logger.Warn
    if os.Getenv("APP_ENV") == "prod" {
        logLevel = logger.Silent
    }

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logLevel),
    })

    if err != nil {
        log.Fatal("Cannot connect to database with given DSN")
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
