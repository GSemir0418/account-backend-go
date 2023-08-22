package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

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
	Email     string `gorm:"uniqueIndex"`
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Item struct {
	ID         int
	UserID     int
	HappenedAt time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
type Tag struct {
	ID int
}

var models = []any{&User{}, &Item{}, &Tag{}}

func CreateTables() {
	for _, model := range models {
		err := DB.Migrator().CreateTable(model)
		if err != nil {
			log.Println(err)
		}
	}
	log.Println("Successfully create tables")
}

func Migrate() {
	// 自动迁移
	// 可以通过修改 User 结构体的内容以及标签来修改字段属性
	// 只能新增，不能删除
	// 为了安全，特意设计的
	DB.AutoMigrate(models...)
}

func Crud() {
}
