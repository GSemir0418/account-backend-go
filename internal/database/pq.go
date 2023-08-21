package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "gsemir"
	password = "gsemir"
	dbname   = "go_account_dev"
)

func Connect() {
	// 构建连接字符串
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	DB = db
}

func Close() {
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalln(err)
	}
	sqlDB.Close()
}

type User struct {
	ID        int
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func CreateTables() {
	err := DB.Migrator().CreateTable(&User{})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Successfully create user table")
}

func Migrate() {
}

func Crud() {
}
