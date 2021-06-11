package database

import (
	"fmt"

	"github.com/leopurba/go-article/config"
	"github.com/leopurba/go-article/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Client interface {
	Conn() *gorm.DB
	Close() error
}

func NewClient() (*gorm.DB, error) {
	dbHost := config.Cfg().PostgresHost
	dbPort := config.Cfg().PostgresPort
	dbUser := config.Cfg().PostgresUser
	dbPass := config.Cfg().PostgresPassword
	dbName := config.Cfg().PostgresDatabase
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Jakarta", dbHost, dbPort, dbUser, dbPass, dbName)

	dbConn, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		Logger: *logger.NewLog(),
	})
	if err != nil {
		return dbConn, err
	}

	db, err := dbConn.DB()
	if err != nil {
		return dbConn, err
	}

	err = db.Ping()
	if err != nil {
		return dbConn, err
	}
	return dbConn, nil
}
