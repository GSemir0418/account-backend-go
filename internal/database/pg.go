package database

import (
	"database/sql"
	"fmt"

	"log"

	// 引入pg的驱动，前面加_表示不直接访问这个包
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "gsemir"
	password = "gsemir"
	dbname   = "go_account_dev"
)

func PgConnect() {
	// 构建连接字符串
	connectStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", connectStr)
	if err != nil {
		log.Fatal(err)
	}
	// 及时在不同文件下，包名相同的变量可以访问
	DB = db
	// 两个步骤共享同一个 err
	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully connected to db")
}

func PgClose() {
	DB.Close()
	log.Println("Closed db")
}

func PgCreateTables() {
	// 创建 Users 表
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(100) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Successfully created user table")
}
