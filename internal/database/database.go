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
	// 创建一个 User
	user := User{Email: "test1@qq.com"}
	// crud 的返回值都是一个事务
	tx := DB.Create(&user)
	log.Println(tx.RowsAffected)
	// 不会将新user返回，而是直接修改 user 实例
	log.Println(user)

	// 查询
	u2 := User{}
	_ = DB.Find(&u2, 1)

	// 更新
	u2.Phone = "123456789"
	tx = DB.Save(&u2)
	if tx.Error != nil {
		log.Println(tx.Error)
	} else {
		log.Println(tx.RowsAffected)
		log.Println(u2)
	}

	// 分页排序查询
	users := []User{}
	// 指定表查询（不指定的话就查询Find中传入的结构体实例对应的表）
	// DB.Model(&User{})
	DB.Offset(0).Limit(10).Order("created_at asc, id desc").Find(&users)
	log.Println(users)

	// 删除
	u := User{ID: 1}
	tx = DB.Delete(&u)
	if tx.Error != nil {
		log.Println(tx.Error)
	} else {
		log.Println(tx.RowsAffected)
	}

	// 裸 sql 查询
	tx = DB.Raw("SELECT * FROM users WHERE id = ?", 4).Scan(&u)
	if tx.Error != nil {
		log.Println(tx.Error)
	} else {
		log.Println(u)
	}

	// 本质上 orm 就是帮我们在写 sql 时分段
	// 不同的orm所对应的api也不同，确实需要一定的学习成本，尤其是当需求比较复杂时，还不如手写sql
	// 并不能够提供类型安全的保障，跟裸写sql大差不差
	// gorm 类型安全插件 gen，官网有，可以酌情使用
}
