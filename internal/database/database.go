package database

import (
	queries "account/config/sqlc"
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/exec"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// 数据库全局变量
var DB *sql.DB

// 数据库上下文
var DBCtx = context.Background()

const (
	host     = "localhost"
	port     = 5432
	user     = "gsemir"
	password = "gsemir"
	dbname   = "go_account_dev"
)

func Connect() {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalln(err)
	}
	DB = db
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
}

func Close() {
}

func MigrateCreate(filename string) {
	// migration 库没有提供创建迁移文件的方法
	// 使用 go 执行语句即可，把命令拆成字符串
	cmd := exec.Command("migrate", "create", "-ext", "sql", "-dir", "config/migrations", "-seq", filename)
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

func MigrateUp() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	m, err := migrate.New(
		fmt.Sprintf("file://%s/config/migrations", dir),
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			user, password, host, port, dbname,
		),
	)
	if err != nil {
		log.Fatalln(err)
	}
	err = m.Up() // 会直接同步所有更新
	if err != nil {
		log.Fatalln(err)
	}

}

func MigrateDown() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	m, err := migrate.New(
		fmt.Sprintf("file://%s/config/migrations", dir),
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			user, password, host, port, dbname,
		),
	)
	if err != nil {
		log.Fatalln(err)
	}
	err = m.Steps(-1) // 默认只回退一步
	if err != nil {
		log.Fatalln(err)
	}
}
func Crud() {
	// queries 来自 sqlc 自动生成的包
	q := queries.New(DB)
	// 随机数字，因为 email 字段是唯一的
	id := rand.Int()

	// 创建 User 返回 User 结构体实例
	u, err := q.CreateUser(DBCtx, fmt.Sprintf("%d@qq.com", id))
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Success create", u)

	// 更新 User
	err = q.UpdateUser(DBCtx, queries.UpdateUserParams{
		ID:    u.ID,
		Email: fmt.Sprintf("%d@qq.com", id),
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Success update")

	// 查询 Users
	users, err := q.ListUsers(DBCtx, queries.ListUsersParams{
		Offset: 0,
		Limit:  10,
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Success Query", users)
}
