package db

import (
	"fmt"
	"github.com/zzz4zzz/go-playground/db/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

type DatabaseConfiguration struct {
	Host            string
	Port            int
	User            string
	Password        string
	DatabaseName    string
	ApplicationName string
}

func (c *DatabaseConfiguration) String() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", c.Host, c.User, c.Password, c.DatabaseName, c.Port)
}

func InitDatabase(config DatabaseConfiguration) *gorm.DB {
	dsn := config.String()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{AllowGlobalUpdate: false, SkipDefaultTransaction: true})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.AutoMigrate(&models.User{}, &models.BankAccount{})
	if err != nil {
		panic("Database migration failed")
	}
	DB = db
	return db
}
