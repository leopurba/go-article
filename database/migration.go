package database

import (
	"github.com/leopurba/go-article/models"

	"gorm.io/gorm"
)

type Migration interface {
	Migration() error
}
type myDB struct {
	DB *gorm.DB
}

func NewMigration(db *gorm.DB) Migration {
	return &myDB{
		DB: db,
	}
}

func (myDB *myDB) Migration() (err error) {
	err = myDB.DB.AutoMigrate(&models.Article{})
	return
}
