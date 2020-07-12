package mysql

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"go/common/log"
	"time"
)

type Config struct {
	// demo:password@tcp(127.0.0.1:3607)/demo?charset=utf-8
	Dsn         string
	MaxConn     int
	MaxIdleConn int
	MaxLifeTime int
}

func NewDB(config Config) *sql.DB {
	db, _ := sql.Open("mysql", config.Dsn)
	db.SetMaxIdleConns(config.MaxIdleConn)
	db.SetConnMaxLifetime(time.Duration(config.MaxLifeTime) * time.Minute)
	db.SetMaxOpenConns(config.MaxConn)
	if err := db.Ping(); err != nil {
		log.Errorf("open mysql error.error is %v", err)
		panic("open mysql error")
	}
	return db
}
