package orm

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/adhocore/urlsh/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var conn *gorm.DB

// pgConnect connects to postgres db
// It returns gorm DB instance.
func pgConnect() *gorm.DB {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("Database configuration DSN missing, Pass in DATABASE_URL env")
	}

	parse, _ := url.Parse(dsn)
	dbname := strings.Trim(parse.Path, "/")
	pass, _ := parse.User.Password()
	host, port, _ := net.SplitHostPort(parse.Host)

	logLevel, env := logger.Warn, os.Getenv("APP_ENV")
	if env == "prod" {
		logLevel = logger.Silent
	} else if env == "test" {
		dbname = "test_" + dbname
	}

	dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s database=%s", host, port, parse.User.Username(), pass, dbname)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		log.Fatalf("database error: %v", err)
	}

	_ = db.AutoMigrate(&model.Keyword{}, &model.URL{})

	return db
}

var once sync.Once

// Connection gets active db connection ensuring connection is made only once.
// It returns gorm DB instance.
func Connection() *gorm.DB {
	once.Do(func() {
		conn = pgConnect()
	})

	return conn
}
